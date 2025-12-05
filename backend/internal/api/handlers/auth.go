package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/joaosantos/jlpt5/internal/api/dto"
	"github.com/joaosantos/jlpt5/internal/domain/services"
	"github.com/joaosantos/jlpt5/internal/utils"
	pkgErrors "github.com/joaosantos/jlpt5/pkg/errors"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService *services.AuthService
	logger      *utils.Logger
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *services.AuthService, logger *utils.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, pkgErrors.BadRequest("Method not allowed"))
		return
	}

	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, pkgErrors.BadRequest("Invalid request body"))
		return
	}

	authResp, err := h.authService.Register(r.Context(), services.RegisterRequest{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		h.sendError(w, err)
		return
	}

	h.sendSuccess(w, http.StatusCreated, h.toAuthResponse(authResp))
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, pkgErrors.BadRequest("Method not allowed"))
		return
	}

	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, pkgErrors.BadRequest("Invalid request body"))
		return
	}

	authResp, err := h.authService.Login(r.Context(), services.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		h.sendError(w, err)
		return
	}

	h.sendSuccess(w, http.StatusOK, h.toAuthResponse(authResp))
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, pkgErrors.BadRequest("Method not allowed"))
		return
	}

	var req dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, pkgErrors.BadRequest("Invalid request body"))
		return
	}

	authResp, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		h.sendError(w, err)
		return
	}

	h.sendSuccess(w, http.StatusOK, h.toAuthResponse(authResp))
}

// toAuthResponse converts service AuthResponse to DTO AuthResponse
func (h *AuthHandler) toAuthResponse(authResp *services.AuthResponse) dto.AuthResponse {
	var lastLogin *string
	if authResp.User.LastLoginAt != nil {
		formatted := authResp.User.LastLoginAt.Format("2006-01-02T15:04:05Z07:00")
		lastLogin = &formatted
	}

	return dto.AuthResponse{
		User: dto.UserResponse{
			ID:          authResp.User.ID,
			Email:       authResp.User.Email,
			Username:    authResp.User.Username,
			IsActive:    authResp.User.IsActive,
			CreatedAt:   authResp.User.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			LastLoginAt: lastLogin,
		},
		AccessToken:  authResp.AccessToken,
		RefreshToken: authResp.RefreshToken,
	}
}

// sendSuccess sends a successful JSON response
func (h *AuthHandler) sendSuccess(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// sendError sends an error JSON response
func (h *AuthHandler) sendError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	// Check if it's an AppError
	if appErr, ok := err.(*pkgErrors.AppError); ok {
		w.WriteHeader(appErr.StatusCode)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": appErr.Message,
			"code":  appErr.Code,
		})
		return
	}

	// Default to internal server error
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": "Internal server error",
		"code":  "INTERNAL_ERROR",
	})
}
