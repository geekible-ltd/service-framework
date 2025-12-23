package frameworkdto

import "time"

type UserUpdateRequestDTO struct {
	UserID          uint   `json:"user_id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	Role            string `json:"role"`
	IsActive        bool   `json:"is_active"`
	IsEmailVerified bool   `json:"is_email_verified"`
}

type ResetPasswordRequestDTO struct {
	Email string `json:"email"`
}

type ResetPasswordDTO struct {
	ResetToken  string `json:"reset_token"`
	NewPassword string `json:"new_password"`
}

type VerifyEmailDTO struct {
	TenantID uint   `json:"tenant_id"`
	UserID   uint   `json:"user_id"`
	Token    string `json:"token"`
}

type GetUsersResponseDTO struct {
	UserID    uint      `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetUserRoles struct {
	Role string `json:"role"`
}
