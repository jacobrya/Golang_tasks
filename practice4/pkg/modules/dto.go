package modules

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type UserResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Age       int    `json:"age"`
	CreatedAt string `json:"created_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
	ID      int    `json:"id,omitempty"`
}

type DeleteResponse struct {
	Message      string `json:"message"`
	RowsAffected int64  `json:"rows_affected"`
}

