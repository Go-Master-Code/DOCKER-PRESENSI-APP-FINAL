package dto

type UserResponse struct {
	ID int `json:"id" `
	// Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	RoleID   int    `json:"role_id"`
	RoleNama string `json:"role_nama"`
}

type CreateUserRequest struct {
	// Email    string `json:"email" binding:"required,email,max=100"` // wajib format email
	Username string `json:"username" binding:"required,max=50"`
	Password string `json:"password" binding:"required,max=50"`
	RoleID   int    `json:"role_id" binding:"required"`
}

type UpdateUserRequest struct {
	// Email    *string `json:"email" binding:"omitempty,email"`
	Username *string `json:"username" binding:"omitempty"`
	Password *string `json:"password" binding:"omitempty"`
	RoleID   *int    `json:"role_id" binding:"omitempty"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,max=50"`
	Password string `json:"password" binding:"required,max=50"`
}
