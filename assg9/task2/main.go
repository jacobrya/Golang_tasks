package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type DBStore struct {
	db *sql.DB
}

func NewDBStore(connStr string) (*DBStore, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS idempotency_keys (
			key         TEXT PRIMARY KEY,
			status      TEXT NOT NULL DEFAULT 'processing',
			status_code INT,
			body        BYTEA,
			created_at  TIMESTAMPTZ DEFAULT NOW()
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &DBStore{db: db}, nil
}

func (s *DBStore) StartProcessing(key string) (bool, error) {
	_, err := s.db.Exec(`
		INSERT INTO idempotency_keys (key, status)
		VALUES ($1, 'processing')
		ON CONFLICT (key) DO NOTHING
	`, key)
	if err != nil {
		return false, err
	}

	var status string
	err = s.db.QueryRow(`SELECT status FROM idempotency_keys WHERE key = $1`, key).Scan(&status)
	if err != nil {
		return false, err
	}

	return status == "processing", nil
}

type CachedResponse struct {
	StatusCode int
	Body       []byte
	Status     string
}

func (s *DBStore) Get(key string) (*CachedResponse, bool) {
	var status string
	var statusCode sql.NullInt64
	var body []byte

	err := s.db.QueryRow(`
		SELECT status, status_code, body FROM idempotency_keys WHERE key = $1
	`, key).Scan(&status, &statusCode, &body)

	if err == sql.ErrNoRows {
		return nil, false
	}
	if err != nil {
		return nil, false
	}

	return &CachedResponse{
		Status:     status,
		StatusCode: int(statusCode.Int64),
		Body:       body,
	}, true
}

func (s *DBStore) Finish(key string, statusCode int, body []byte) error {
	_, err := s.db.Exec(`
		UPDATE idempotency_keys
		SET status = 'completed', status_code = $2, body = $3
		WHERE key = $1
	`, key, statusCode, body)
	return err
}

func IdempotencyMiddleware(store *DBStore, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Idempotency-Key")
		if key == "" {
			http.Error(w, "Idempotency-Key header required", http.StatusBadRequest)
			return
		}

		if cached, exists := store.Get(key); exists {
			if cached.Status == "completed" {
				fmt.Println("[Middleware] Key already completed, returning cached response")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(cached.StatusCode)
				w.Write(cached.Body)
				return
			}
			fmt.Println("[Middleware] Duplicate request in progress → 409")
			http.Error(w, "Duplicate request in progress", http.StatusConflict)
			return
		}

		ok, err := store.StartProcessing(key)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if !ok {
			if cached, exists := store.Get(key); exists && cached.Status == "completed" {
				fmt.Println("[Middleware] Just completed by another request, returning cache")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(cached.StatusCode)
				w.Write(cached.Body)
			} else {
				fmt.Println("[Middleware] Race: duplicate in progress → 409")
				http.Error(w, "Duplicate request in progress", http.StatusConflict)
			}
			return
		}

		recorder := httptest.NewRecorder()
		next.ServeHTTP(recorder, r)

		if err := store.Finish(key, recorder.Code, recorder.Body.Bytes()); err != nil {
			fmt.Printf("[Middleware] Failed to save result: %v\n", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(recorder.Code)
		w.Write(recorder.Body.Bytes())
	}
}

func PaymentHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[Handler] Processing started...")
	time.Sleep(2 * time.Second)

	response := map[string]interface{}{
		"status":         "paid",
		"amount":         1000,
		"transaction_id": uuid.New().String(),
	}
	body, _ := json.Marshal(response)

	fmt.Println("[Handler] Processing completed")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", ""),
		getEnv("DB_NAME", "postgres"),
	)

	store, err := NewDBStore(connStr)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to PostgreSQL: %v", err))
	}
	fmt.Println("Connected to PostgreSQL")

	handler := IdempotencyMiddleware(store, PaymentHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	fmt.Println("\n=== Simulating Double-Click Attack ===")
	key := uuid.New().String()
	fmt.Printf("Using Idempotency-Key: %s\n\n", key)

	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			req, _ := http.NewRequest("POST", server.URL, nil)
			req.Header.Set("Idempotency-Key", key)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Printf("[Request %d] Error: %v\n", id, err)
				return
			}
			defer resp.Body.Close()

			fmt.Printf("[Request %d] Status: %d\n", id, resp.StatusCode)
		}(i)
	}

	wg.Wait()

	fmt.Println("\n=== Final request after completion ===")
	req, _ := http.NewRequest("POST", server.URL, nil)
	req.Header.Set("Idempotency-Key", key)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	fmt.Printf("Final Request: Status %d\n", resp.StatusCode)
}
