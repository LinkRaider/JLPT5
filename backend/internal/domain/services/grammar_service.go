package services

import (
	"context"
	"time"

	"github.com/joaosantos/jlpt5/internal/domain/models"
	"github.com/joaosantos/jlpt5/internal/domain/repository"
	"github.com/joaosantos/jlpt5/internal/utils"
	pkgErrors "github.com/joaosantos/jlpt5/pkg/errors"
)

// GrammarService handles grammar business logic
type GrammarService struct {
	grammarRepo repository.GrammarRepository
	logger      *utils.Logger
}

// NewGrammarService creates a new grammar service
func NewGrammarService(grammarRepo repository.GrammarRepository, logger *utils.Logger) *GrammarService {
	return &GrammarService{
		grammarRepo: grammarRepo,
		logger:      logger,
	}
}

// GetLessonsList retrieves grammar lessons with optional filtering
func (s *GrammarService) GetLessonsList(ctx context.Context, userID int, jlptLevel *int, page, pageSize int) ([]models.GrammarLessonWithExamples, int, error) {
	offset := (page - 1) * pageSize

	lessons, err := s.grammarRepo.GetUserLessonsList(ctx, userID, jlptLevel, pageSize, offset)
	if err != nil {
		s.logger.Error("Failed to get grammar lessons", utils.WithContext("error", err.Error()))
		return nil, 0, pkgErrors.Internal("Failed to retrieve grammar lessons", err)
	}

	total, err := s.grammarRepo.CountLessons(ctx, jlptLevel)
	if err != nil {
		s.logger.Error("Failed to count grammar lessons", utils.WithContext("error", err.Error()))
		return nil, 0, pkgErrors.Internal("Failed to count grammar lessons", err)
	}

	return lessons, total, nil
}

// GetLessonByID retrieves a specific grammar lesson with examples
func (s *GrammarService) GetLessonByID(ctx context.Context, userID, lessonID int) (*models.GrammarLessonWithExamples, error) {
	lesson, err := s.grammarRepo.GetLessonByID(ctx, lessonID)
	if err != nil {
		return nil, err
	}

	// Try to get user progress
	progress, err := s.grammarRepo.GetUserProgress(ctx, userID, lessonID)
	if err == nil {
		lesson.Progress = progress
	}

	return lesson, nil
}

// MarkAsCompleted marks a grammar lesson as completed
func (s *GrammarService) MarkAsCompleted(ctx context.Context, userID, lessonID int, notes *string) (*models.UserGrammarProgress, error) {
	// Verify lesson exists
	_, err := s.grammarRepo.GetLessonByID(ctx, lessonID)
	if err != nil {
		return nil, err
	}

	// Get or create progress
	progress, err := s.grammarRepo.GetUserProgress(ctx, userID, lessonID)
	if err != nil {
		if appErr, ok := err.(*pkgErrors.AppError); ok && appErr.Code == pkgErrors.ErrCodeNotFound {
			// Create new progress
			now := time.Now()
			progress = &models.UserGrammarProgress{
				UserID:          userID,
				GrammarLessonID: lessonID,
				Completed:       true,
				CompletedAt:     &now,
				Notes:           notes,
			}

			if createErr := s.grammarRepo.CreateUserProgress(ctx, progress); createErr != nil {
				s.logger.Error("Failed to create progress", utils.WithContext("error", createErr.Error()))
				return nil, pkgErrors.Internal("Failed to create progress", createErr)
			}

			s.logger.Info("Grammar lesson marked as completed", utils.WithContext("user_id", userID, "lesson_id", lessonID))
			return progress, nil
		}
		return nil, err
	}

	// Update existing progress
	now := time.Now()
	progress.Completed = true
	progress.CompletedAt = &now
	if notes != nil {
		progress.Notes = notes
	}

	if err := s.grammarRepo.UpdateUserProgress(ctx, progress); err != nil {
		s.logger.Error("Failed to update progress", utils.WithContext("error", err.Error()))
		return nil, pkgErrors.Internal("Failed to update progress", err)
	}

	s.logger.Info("Grammar lesson updated", utils.WithContext("user_id", userID, "lesson_id", lessonID))
	return progress, nil
}
