package repository

import (
	"context"

	"github.com/joaosantos/jlpt5/internal/domain/models"
)

// VocabularyRepository defines the interface for vocabulary data access
type VocabularyRepository interface {
	// GetAll retrieves all vocabulary items with optional filtering
	GetAll(ctx context.Context, jlptLevel *int, limit, offset int) ([]models.Vocabulary, error)

	// GetByID retrieves a vocabulary item by ID
	GetByID(ctx context.Context, id int) (*models.Vocabulary, error)

	// GetDueForReview retrieves vocabulary items due for review for a user
	GetDueForReview(ctx context.Context, userID int, limit int) ([]models.VocabularyWithProgress, error)

	// GetUserProgress retrieves user's progress for a vocabulary item
	GetUserProgress(ctx context.Context, userID, vocabularyID int) (*models.UserVocabularyProgress, error)

	// CreateUserProgress creates initial progress for a user-vocabulary pair
	CreateUserProgress(ctx context.Context, progress *models.UserVocabularyProgress) error

	// UpdateUserProgress updates user's progress for a vocabulary item
	UpdateUserProgress(ctx context.Context, progress *models.UserVocabularyProgress) error

	// GetUserVocabularyList retrieves vocabulary with user progress
	GetUserVocabularyList(ctx context.Context, userID int, jlptLevel *int, limit, offset int) ([]models.VocabularyWithProgress, error)

	// Count returns total vocabulary count
	Count(ctx context.Context, jlptLevel *int) (int, error)
}
