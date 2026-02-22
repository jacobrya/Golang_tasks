package users

import (
	"database/sql"
	"errors"
	"fmt"
	"prac4/internal/repository/_postgres"
	"prac4/pkg/modules"
	"time"
)

type Repository struct {
	db               *_postgres.Dialect
	executionTimeout time.Duration
}

func NewUserRepository(db *_postgres.Dialect) *Repository {
	return &Repository{
		db:               db,
		executionTimeout: time.Second * 5,
	}
}

func (r *Repository) GetUsers() ([]modules.User, error) {
	var users []modules.User
	err := r.db.DB.Select(&users, "SELECT id, name, email, age, created_at FROM users")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}
	fmt.Println(users)
	return users, nil
}

// CreateUser creates a new user and returns the newly generated ID
func (r *Repository) CreateUser(user *modules.User) (int, error) {
	if user == nil {
		return 0, errors.New("user cannot be nil")
	}

	query := `
		INSERT INTO users (name, email, age, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var id int
	err := r.db.DB.QueryRow(query, user.Name, user.Email, user.Age, time.Now()).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}

// UpdateUser updates an existing user and returns an error if the user doesn't exist
func (r *Repository) UpdateUser(user *modules.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	if user.ID == 0 {
		return errors.New("user ID is required for update")
	}

	query := `
		UPDATE users
		SET name = $1, email = $2, age = $3
		WHERE id = $4
	`

	result, err := r.db.DB.Exec(query, user.Name, user.Email, user.Age, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d does not exist", user.ID)
	}

	return nil
}

// GetUserByID fetches a single user by ID, returns nil and error if not found
func (r *Repository) GetUserByID(id int) (*modules.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user ID")
	}

	var user modules.User
	query := "SELECT id, name, email, age, created_at FROM users WHERE id = $1"

	err := r.db.DB.Get(&user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with ID %d does not exist", id)
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	return &user, nil
}

// DeleteUserByID deletes a user by ID and returns the number of rows affected
func (r *Repository) DeleteUserByID(id int) (int64, error) {
	if id <= 0 {
		return 0, errors.New("invalid user ID")
	}

	query := "DELETE FROM users WHERE id = $1"

	result, err := r.db.DB.Exec(query, id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return 0, fmt.Errorf("user with ID %d does not exist", id)
	}

	return rowsAffected, nil
}
