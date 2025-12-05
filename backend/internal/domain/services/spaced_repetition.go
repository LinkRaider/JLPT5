package services

import (
	"math"
	"time"

	"github.com/joaosantos/jlpt5/internal/domain/models"
)

// ReviewQuality represents how well the user remembered the item
type ReviewQuality int

const (
	// ReviewQualityBlackout - complete blackout (0)
	ReviewQualityBlackout ReviewQuality = 0
	// ReviewQualityIncorrect - incorrect response (1)
	ReviewQualityIncorrect ReviewQuality = 1
	// ReviewQualityIncorrectEasy - incorrect but remembered (2)
	ReviewQualityIncorrectEasy ReviewQuality = 2
	// ReviewQualityCorrectHard - correct with difficulty (3)
	ReviewQualityCorrectHard ReviewQuality = 3
	// ReviewQualityCorrectEasy - correct with hesitation (4)
	ReviewQualityCorrectEasy ReviewQuality = 4
	// ReviewQualityPerfect - perfect response (5)
	ReviewQualityPerfect ReviewQuality = 5
)

// SpacedRepetitionService implements the SM-2 algorithm
type SpacedRepetitionService struct{}

// NewSpacedRepetitionService creates a new spaced repetition service
func NewSpacedRepetitionService() *SpacedRepetitionService {
	return &SpacedRepetitionService{}
}

// CalculateNextReview calculates the next review date using SM-2 algorithm
// Based on SuperMemo 2 algorithm: https://www.supermemo.com/en/archives1990-2015/english/ol/sm2
func (s *SpacedRepetitionService) CalculateNextReview(
	progress *models.UserVocabularyProgress,
	quality ReviewQuality,
) (*models.UserVocabularyProgress, error) {
	// Clone the progress to avoid modifying the original
	newProgress := *progress
	now := time.Now()

	// Update counters
	newProgress.TotalReviews++
	newProgress.LastReviewedAt = &now

	// Quality >= 3 is considered correct
	isCorrect := quality >= ReviewQualityCorrectHard

	if isCorrect {
		newProgress.CorrectReviews++
	}

	// SM-2 Algorithm implementation
	if quality < ReviewQualityCorrectHard {
		// If quality < 3 (incorrect), reset repetitions and interval
		newProgress.Repetitions = 0
		newProgress.Interval = 1
	} else {
		// If quality >= 3 (correct)
		if newProgress.Repetitions == 0 {
			newProgress.Interval = 1
		} else if newProgress.Repetitions == 1 {
			newProgress.Interval = 6
		} else {
			// For repetitions >= 2: I(n) = I(n-1) * EF
			newProgress.Interval = int(math.Round(float64(newProgress.Interval) * newProgress.EaseFactor))
		}

		newProgress.Repetitions++
	}

	// Update ease factor
	// EF' = EF + (0.1 - (5 - q) * (0.08 + (5 - q) * 0.02))
	// Where q is the quality of response (0-5)
	q := float64(quality)
	newProgress.EaseFactor = newProgress.EaseFactor + (0.1 - (5-q)*(0.08+(5-q)*0.02))

	// Ensure ease factor doesn't go below 1.3 (minimum per SM-2 algorithm)
	if newProgress.EaseFactor < 1.3 {
		newProgress.EaseFactor = 1.3
	}

	// Calculate next review date
	newProgress.NextReviewDate = now.AddDate(0, 0, newProgress.Interval)

	return &newProgress, nil
}

// InitializeProgress creates initial progress for a new vocabulary item
func (s *SpacedRepetitionService) InitializeProgress(userID, vocabularyID int) *models.UserVocabularyProgress {
	now := time.Now()
	return &models.UserVocabularyProgress{
		UserID:         userID,
		VocabularyID:   vocabularyID,
		EaseFactor:     2.5, // Default starting ease factor per SM-2
		Interval:       1,   // Start with 1 day interval
		Repetitions:    0,   // No repetitions yet
		NextReviewDate: now, // Available for review immediately
		TotalReviews:   0,
		CorrectReviews: 0,
	}
}

// GetQualityFromBoolean converts a simple correct/incorrect to ReviewQuality
func (s *SpacedRepetitionService) GetQualityFromBoolean(isCorrect bool) ReviewQuality {
	if isCorrect {
		return ReviewQualityCorrectEasy // Default to quality 4 for correct
	}
	return ReviewQualityIncorrect // Default to quality 1 for incorrect
}

// GetReviewStats calculates review statistics for display
func (s *SpacedRepetitionService) GetReviewStats(progress *models.UserVocabularyProgress) map[string]interface{} {
	var successRate float64
	if progress.TotalReviews > 0 {
		successRate = float64(progress.CorrectReviews) / float64(progress.TotalReviews) * 100
	}

	var daysSinceLastReview *int
	if progress.LastReviewedAt != nil {
		days := int(time.Since(*progress.LastReviewedAt).Hours() / 24)
		daysSinceLastReview = &days
	}

	daysUntilNextReview := int(time.Until(progress.NextReviewDate).Hours() / 24)
	if daysUntilNextReview < 0 {
		daysUntilNextReview = 0
	}

	return map[string]interface{}{
		"success_rate":           successRate,
		"total_reviews":          progress.TotalReviews,
		"correct_reviews":        progress.CorrectReviews,
		"current_interval_days":  progress.Interval,
		"repetitions":            progress.Repetitions,
		"ease_factor":            progress.EaseFactor,
		"days_since_last_review": daysSinceLastReview,
		"days_until_next_review": daysUntilNextReview,
		"is_due":                 time.Now().After(progress.NextReviewDate),
	}
}
