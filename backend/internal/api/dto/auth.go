package dto

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginRequest represents a user login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RefreshTokenRequest represents a token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// AuthResponse represents an authentication response
type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
}

// UserResponse represents a user in API responses
type UserResponse struct {
	ID          int     `json:"id"`
	Email       string  `json:"email"`
	Username    string  `json:"username"`
	IsActive    bool    `json:"is_active"`
	CreatedAt   string  `json:"created_at"`
	LastLoginAt *string `json:"last_login_at,omitempty"`
}
