package repository

import (
	"context"

	"github.com/joaosantos/jlpt5/internal/domain/models"
)

// GrammarRepository defines the interface for grammar data access
type GrammarRepository interface {
	// GetAllLessons retrieves all grammar lessons with optional filtering
	GetAllLessons(ctx context.Context, jlptLevel *int, limit, offset int) ([]models.GrammarLesson, error)

	// GetLessonByID retrieves a grammar lesson by ID with examples
	GetLessonByID(ctx context.Context, lessonID int) (*models.GrammarLessonWithExamples, error)

	// GetUserProgress retrieves user's progress for a lesson
	GetUserProgress(ctx context.Context, userID, lessonID int) (*models.UserGrammarProgress, error)

	// CreateUserProgress creates progress for a user-lesson pair
	CreateUserProgress(ctx context.Context, progress *models.UserGrammarProgress) error

	// UpdateUserProgress updates user's progress for a lesson
	UpdateUserProgress(ctx context.Context, progress *models.UserGrammarProgress) error

	// GetUserLessonsList retrieves lessons with user progress
	GetUserLessonsList(ctx context.Context, userID int, jlptLevel *int, limit, offset int) ([]models.GrammarLessonWithExamples, error)

	// CountLessons returns total lesson count
	CountLessons(ctx context.Context, jlptLevel *int) (int, error)
}
