package repository

import (
	"context"

	"github.com/joaosantos/jlpt5/internal/domain/models"
)

// QuizRepository defines the interface for quiz data access
type QuizRepository interface {
	// GetAllQuizzes retrieves all quizzes with optional filtering
	GetAllQuizzes(ctx context.Context, jlptLevel *int, limit, offset int) ([]models.Quiz, error)

	// GetQuizByID retrieves a quiz by ID
	GetQuizByID(ctx context.Context, quizID int) (*models.Quiz, error)

	// GetQuizQuestions retrieves all questions for a quiz
	GetQuizQuestions(ctx context.Context, quizID int) ([]models.QuizQuestion, error)

	// CreateQuizSession creates a new quiz session
	CreateQuizSession(ctx context.Context, session *models.QuizSession) error

	// UpdateQuizSession updates a quiz session
	UpdateQuizSession(ctx context.Context, session *models.QuizSession) error

	// GetQuizSession retrieves a quiz session by ID
	GetQuizSession(ctx context.Context, sessionID int) (*models.QuizSession, error)

	// GetUserQuizSessions retrieves all quiz sessions for a user
	GetUserQuizSessions(ctx context.Context, userID int, limit, offset int) ([]models.QuizSession, error)

	// CreateQuizAnswer creates a quiz answer
	CreateQuizAnswer(ctx context.Context, answer *models.QuizAnswer) error

	// GetSessionAnswers retrieves all answers for a session
	GetSessionAnswers(ctx context.Context, sessionID int) ([]models.QuizAnswer, error)

	// CountQuizzes returns total quiz count
	CountQuizzes(ctx context.Context, jlptLevel *int) (int, error)
}
