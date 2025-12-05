package repository

import (
	"context"

	"github.com/joaosantos/jlpt5/internal/domain/models"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *models.User) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id int) (*models.User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*models.User, error)

	// GetByUsername retrieves a user by username
	GetByUsername(ctx context.Context, username string) (*models.User, error)

	// Update updates a user
	Update(ctx context.Context, user *models.User) error

	// UpdateLastLogin updates the user's last login timestamp
	UpdateLastLogin(ctx context.Context, userID int) error

	// GetStatistics retrieves user statistics
	GetStatistics(ctx context.Context, userID int) (*models.UserStatistics, error)

	// CreateStatistics creates initial statistics for a user
	CreateStatistics(ctx context.Context, userID int) error

	// UpdateStatistics updates user statistics
	UpdateStatistics(ctx context.Context, stats *models.UserStatistics) error
}
