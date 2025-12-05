package services

import (
	"context"

	"github.com/joaosantos/jlpt5/internal/domain/models"
	"github.com/joaosantos/jlpt5/internal/domain/repository"
	"github.com/joaosantos/jlpt5/internal/utils"
	pkgErrors "github.com/joaosantos/jlpt5/pkg/errors"
)

// VocabularyService handles vocabulary business logic
type VocabularyService struct {
	vocabRepo repository.VocabularyRepository
	srService *SpacedRepetitionService
	logger    *utils.Logger
}

// NewVocabularyService creates a new vocabulary service
func NewVocabularyService(
	vocabRepo repository.VocabularyRepository,
	srService *SpacedRepetitionService,
	logger *utils.Logger,
) *VocabularyService {
	return &VocabularyService{
		vocabRepo: vocabRepo,
		srService: srService,
		logger:    logger,
	}
}

// GetVocabularyList retrieves vocabulary items with optional filtering
func (s *VocabularyService) GetVocabularyList(ctx context.Context, userID int, jlptLevel *int, page, pageSize int) ([]models.VocabularyWithProgress, int, error) {
	offset := (page - 1) * pageSize

	items, err := s.vocabRepo.GetUserVocabularyList(ctx, userID, jlptLevel, pageSize, offset)
	if err != nil {
		s.logger.Error("Failed to get vocabulary list", utils.WithContext("error", err.Error()))
		return nil, 0, pkgErrors.Internal("Failed to retrieve vocabulary", err)
	}

	total, err := s.vocabRepo.Count(ctx, jlptLevel)
	if err != nil {
		s.logger.Error("Failed to count vocabulary", utils.WithContext("error", err.Error()))
		return nil, 0, pkgErrors.Internal("Failed to count vocabulary", err)
	}

	return items, total, nil
}

// GetVocabularyByID retrieves a specific vocabulary item
func (s *VocabularyService) GetVocabularyByID(ctx context.Context, userID, vocabularyID int) (*models.VocabularyWithProgress, error) {
	vocab, err := s.vocabRepo.GetByID(ctx, vocabularyID)
	if err != nil {
		return nil, err
	}

	result := &models.VocabularyWithProgress{
		Vocabulary: *vocab,
	}

	// Try to get user progress
	progress, err := s.vocabRepo.GetUserProgress(ctx, userID, vocabularyID)
	if err == nil {
		result.Progress = progress
	} else {
		// Check if it's not a "not found" error, then log it
		if _, ok := err.(*pkgErrors.AppError); !ok || err.(*pkgErrors.AppError).Code != pkgErrors.ErrCodeNotFound {
			s.logger.Error("Failed to get user progress", utils.WithContext("error", err.Error()))
		}
	}

	return result, nil
}

// GetDueVocabulary retrieves vocabulary items due for review
func (s *VocabularyService) GetDueVocabulary(ctx context.Context, userID int, limit int) ([]models.VocabularyWithProgress, error) {
	items, err := s.vocabRepo.GetDueForReview(ctx, userID, limit)
	if err != nil {
		s.logger.Error("Failed to get due vocabulary", utils.WithContext("error", err.Error()))
		return nil, pkgErrors.Internal("Failed to retrieve due vocabulary", err)
	}

	return items, nil
}

// SubmitReview processes a review submission
func (s *VocabularyService) SubmitReview(ctx context.Context, userID, vocabularyID int, isCorrect bool) (*models.UserVocabularyProgress, error) {
	// Get or create progress
	progress, err := s.vocabRepo.GetUserProgress(ctx, userID, vocabularyID)
	if err != nil {
		// If progress doesn't exist, create it
		if appErr, ok := err.(*pkgErrors.AppError); ok && appErr.Code == pkgErrors.ErrCodeNotFound {
			progress = s.srService.InitializeProgress(userID, vocabularyID)
			if createErr := s.vocabRepo.CreateUserProgress(ctx, progress); createErr != nil {
				s.logger.Error("Failed to create progress", utils.WithContext("error", createErr.Error()))
				return nil, pkgErrors.Internal("Failed to create progress", createErr)
			}
			// Refetch to get the ID
			progress, err = s.vocabRepo.GetUserProgress(ctx, userID, vocabularyID)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	// Calculate new progress using SM-2 algorithm
	quality := s.srService.GetQualityFromBoolean(isCorrect)
	newProgress, err := s.srService.CalculateNextReview(progress, quality)
	if err != nil {
		s.logger.Error("Failed to calculate next review", utils.WithContext("error", err.Error()))
		return nil, pkgErrors.Internal("Failed to calculate review", err)
	}

	// Update progress in database
	if err := s.vocabRepo.UpdateUserProgress(ctx, newProgress); err != nil {
		s.logger.Error("Failed to update progress", utils.WithContext("error", err.Error()))
		return nil, pkgErrors.Internal("Failed to update progress", err)
	}

	s.logger.Info("Review submitted", utils.WithContext(
		"user_id", userID,
		"vocabulary_id", vocabularyID,
		"is_correct", isCorrect,
		"new_interval", newProgress.Interval,
	))

	return newProgress, nil
}

// StartStudying initializes progress for a vocabulary item (marks it as "started")
func (s *VocabularyService) StartStudying(ctx context.Context, userID, vocabularyID int) (*models.UserVocabularyProgress, error) {
	// Check if progress already exists
	_, err := s.vocabRepo.GetUserProgress(ctx, userID, vocabularyID)
	if err == nil {
		return nil, pkgErrors.Conflict("Already studying this vocabulary")
	}

	// Verify vocabulary exists
	_, err = s.vocabRepo.GetByID(ctx, vocabularyID)
	if err != nil {
		return nil, err
	}

	// Create initial progress
	progress := s.srService.InitializeProgress(userID, vocabularyID)
	if err := s.vocabRepo.CreateUserProgress(ctx, progress); err != nil {
		s.logger.Error("Failed to create progress", utils.WithContext("error", err.Error()))
		return nil, pkgErrors.Internal("Failed to start studying", err)
	}

	s.logger.Info("Started studying vocabulary", utils.WithContext(
		"user_id", userID,
		"vocabulary_id", vocabularyID,
	))

	return progress, nil
}

// GetReviewStats gets detailed statistics for a user's vocabulary progress
func (s *VocabularyService) GetReviewStats(ctx context.Context, userID, vocabularyID int) (map[string]interface{}, error) {
	progress, err := s.vocabRepo.GetUserProgress(ctx, userID, vocabularyID)
	if err != nil {
		return nil, err
	}

	return s.srService.GetReviewStats(progress), nil
}
