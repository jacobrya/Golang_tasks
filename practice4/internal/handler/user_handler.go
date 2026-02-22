package handler

import (
	"encoding/json"
	"net/http"
	"prac4/internal/usecase"
	"prac4/pkg/modules"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	usecase *usecase.UserUsecase
}

func NewUserHandler(usecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		usecase: usecase,
	}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.usecase.GetUsers()
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := make([]modules.UserResponse, len(users))
	for i, user := range users {
		response[i] = modules.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Age:       user.Age,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		}
	}

	h.sendJSON(w, http.StatusOK, response)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	user, err := h.usecase.GetUserByID(id)
	if err != nil {
		h.sendError(w, http.StatusNotFound, err.Error())
		return
	}

	response := modules.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Age:       user.Age,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}

	h.sendJSON(w, http.StatusOK, response)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req modules.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" || req.Email == "" {
		h.sendError(w, http.StatusBadRequest, "name and email are required")
		return
	}

	user := &modules.User{
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}

	id, err := h.usecase.CreateUser(user)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := modules.SuccessResponse{
		Message: "User created successfully",
		ID:      id,
	}

	h.sendJSON(w, http.StatusCreated, response)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	var req modules.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user := &modules.User{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}

	err = h.usecase.UpdateUser(user)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			h.sendError(w, http.StatusNotFound, err.Error())
			return
		}
		h.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := modules.SuccessResponse{
		Message: "User updated successfully",
	}

	h.sendJSON(w, http.StatusOK, response)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	rowsAffected, err := h.usecase.DeleteUserByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			h.sendError(w, http.StatusNotFound, err.Error())
			return
		}
		h.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := modules.DeleteResponse{
		Message:      "User deleted successfully",
		RowsAffected: rowsAffected,
	}

	h.sendJSON(w, http.StatusOK, response)
}

func (h *UserHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "ok",
	}
	h.sendJSON(w, http.StatusOK, response)
}

func (h *UserHandler) sendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *UserHandler) sendError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(modules.ErrorResponse{Error: message})
}

