package services

import (
	"context"

	"github.com/joaosantos/jlpt5/internal/domain/models"
	"github.com/joaosantos/jlpt5/internal/domain/repository"
	"github.com/joaosantos/jlpt5/internal/utils"
	pkgErrors "github.com/joaosantos/jlpt5/pkg/errors"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo   repository.UserRepository
	jwtManager *utils.JWTManager
	logger     *utils.Logger
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepository, jwtManager *utils.JWTManager, logger *utils.Logger) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
		logger:     logger,
	}
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Email    string
	Username string
	Password string
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string
	Password string
}

// AuthResponse represents an authentication response
type AuthResponse struct {
	User         *models.User
	AccessToken  string
	RefreshToken string
}

// Register registers a new user
func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	// Validate input
	if req.Email == "" || req.Username == "" || req.Password == "" {
		return nil, pkgErrors.Validation("Email, username, and password are required")
	}

	if len(req.Password) < 8 {
		return nil, pkgErrors.Validation("Password must be at least 8 characters")
	}

	// Hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		s.logger.Error("Failed to hash password", utils.WithContext("error", err.Error()))
		return nil, pkgErrors.Internal("Failed to create user", err)
	}

	// Create user
	user := &models.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: passwordHash,
		IsActive:     true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.Error("Failed to create user", utils.WithContext("error", err.Error(), "email", req.Email))
		return nil, err
	}

	// Create initial statistics for the user
	if err := s.userRepo.CreateStatistics(ctx, user.ID); err != nil {
		s.logger.Error("Failed to create user statistics", utils.WithContext("error", err.Error(), "user_id", user.ID))
		// Not a fatal error, continue
	}

	// Generate tokens
	accessToken, err := s.jwtManager.GenerateToken(user.ID, user.Email, user.Username)
	if err != nil {
		s.logger.Error("Failed to generate access token", utils.WithContext("error", err.Error()))
		return nil, pkgErrors.Internal("Failed to generate token", err)
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Email, user.Username)
	if err != nil {
		s.logger.Error("Failed to generate refresh token", utils.WithContext("error", err.Error()))
		return nil, pkgErrors.Internal("Failed to generate refresh token", err)
	}

	s.logger.Info("User registered successfully", utils.WithContext("user_id", user.ID, "email", user.Email))

	return &AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Login authenticates a user
func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	// Validate input
	if req.Email == "" || req.Password == "" {
		return nil, pkgErrors.Validation("Email and password are required")
	}

	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Warn("Login attempt with invalid email", utils.WithContext("email", req.Email))
		return nil, pkgErrors.InvalidCredentials()
	}

	// Check if user is active
	if !user.IsActive {
		return nil, pkgErrors.Forbidden("Account is inactive")
	}

	// Check password
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		s.logger.Warn("Login attempt with invalid password", utils.WithContext("user_id", user.ID, "email", user.Email))
		return nil, pkgErrors.InvalidCredentials()
	}

	// Update last login
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		s.logger.Error("Failed to update last login", utils.WithContext("error", err.Error(), "user_id", user.ID))
		// Not a fatal error, continue
	}

	// Generate tokens
	accessToken, err := s.jwtManager.GenerateToken(user.ID, user.Email, user.Username)
	if err != nil {
		s.logger.Error("Failed to generate access token", utils.WithContext("error", err.Error()))
		return nil, pkgErrors.Internal("Failed to generate token", err)
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Email, user.Username)
	if err != nil {
		s.logger.Error("Failed to generate refresh token", utils.WithContext("error", err.Error()))
		return nil, pkgErrors.Internal("Failed to generate refresh token", err)
	}

	s.logger.Info("User logged in successfully", utils.WithContext("user_id", user.ID, "email", user.Email))

	return &AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// RefreshToken refreshes an access token using a refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	// Validate refresh token
	claims, err := s.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		return nil, pkgErrors.TokenInvalid()
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	// Generate new tokens
	accessToken, err := s.jwtManager.GenerateToken(user.ID, user.Email, user.Username)
	if err != nil {
		s.logger.Error("Failed to generate access token", utils.WithContext("error", err.Error()))
		return nil, pkgErrors.Internal("Failed to generate token", err)
	}

	newRefreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Email, user.Username)
	if err != nil {
		s.logger.Error("Failed to generate refresh token", utils.WithContext("error", err.Error()))
		return nil, pkgErrors.Internal("Failed to generate refresh token", err)
	}

	s.logger.Info("Token refreshed successfully", utils.WithContext("user_id", user.ID))

	return &AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}
