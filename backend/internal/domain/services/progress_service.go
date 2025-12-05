package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/joaosantos/jlpt5/internal/infrastructure/database"
	"github.com/joaosantos/jlpt5/internal/utils"
)

// ProgressService handles progress and statistics business logic
type ProgressService struct {
	db     *database.DB
	logger *utils.Logger
}

// NewProgressService creates a new progress service
func NewProgressService(db *database.DB, logger *utils.Logger) *ProgressService {
	return &ProgressService{
		db:     db,
		logger: logger,
	}
}

// UserStatistics represents overall user statistics
type UserStatistics struct {
	UserID                 int     `json:"user_id"`
	StudyStreakDays        int     `json:"study_streak_days"`
	LastStudyDate          *string `json:"last_study_date,omitempty"`
	TotalStudyTimeMinutes  int     `json:"total_study_time_minutes"`
	VocabularyLearned      int     `json:"vocabulary_learned"`
	VocabularyDueCount     int     `json:"vocabulary_due_count"`
	GrammarCompleted       int     `json:"grammar_completed"`
	GrammarTotal           int     `json:"grammar_total"`
	QuizzesTaken           int     `json:"quizzes_taken"`
	QuizzesPassed          int     `json:"quizzes_passed"`
	AverageQuizScore       float64 `json:"average_quiz_score"`
}

// GetUserStatistics retrieves comprehensive user statistics
func (s *ProgressService) GetUserStatistics(ctx context.Context, userID int) (*UserStatistics, error) {
	stats := &UserStatistics{UserID: userID}

	// Get basic statistics from user_statistics table
	query := `
		SELECT study_streak_days, last_study_date, total_study_time_minutes,
		       vocabulary_learned, grammar_completed, quizzes_taken, quizzes_passed
		FROM user_statistics
		WHERE user_id = $1
	`

	var lastStudyDate sql.NullTime
	err := s.db.QueryRowContext(ctx, query, userID).Scan(
		&stats.StudyStreakDays, &lastStudyDate, &stats.TotalStudyTimeMinutes,
		&stats.VocabularyLearned, &stats.GrammarCompleted,
		&stats.QuizzesTaken, &stats.QuizzesPassed,
	)

	if err == sql.ErrNoRows {
		// User statistics don't exist yet, return default values
		stats.StudyStreakDays = 0
		stats.TotalStudyTimeMinutes = 0
		stats.VocabularyLearned = 0
		stats.GrammarCompleted = 0
		stats.QuizzesTaken = 0
		stats.QuizzesPassed = 0
	} else if err != nil {
		return nil, fmt.Errorf("error getting user statistics: %w", err)
	}

	if lastStudyDate.Valid {
		dateStr := lastStudyDate.Time.Format("2006-01-02")
		stats.LastStudyDate = &dateStr
	}

	// Get vocabulary due count
	dueCountQuery := `
		SELECT COUNT(*)
		FROM user_vocabulary_progress
		WHERE user_id = $1 AND next_review_date <= CURRENT_TIMESTAMP
	`
	err = s.db.QueryRowContext(ctx, dueCountQuery, userID).Scan(&stats.VocabularyDueCount)
	if err != nil {
		s.logger.Warn("Failed to get vocabulary due count", utils.WithContext("error", err.Error()))
		stats.VocabularyDueCount = 0
	}

	// Get total grammar count
	totalGrammarQuery := `SELECT COUNT(*) FROM grammar_lessons`
	err = s.db.QueryRowContext(ctx, totalGrammarQuery).Scan(&stats.GrammarTotal)
	if err != nil {
		s.logger.Warn("Failed to get total grammar count", utils.WithContext("error", err.Error()))
		stats.GrammarTotal = 0
	}

	// Calculate average quiz score
	if stats.QuizzesTaken > 0 {
		avgScoreQuery := `
			SELECT AVG(percentage)
			FROM quiz_sessions
			WHERE user_id = $1 AND completed_at IS NOT NULL AND percentage IS NOT NULL
		`
		var avgScore sql.NullFloat64
		err = s.db.QueryRowContext(ctx, avgScoreQuery, userID).Scan(&avgScore)
		if err == nil && avgScore.Valid {
			stats.AverageQuizScore = avgScore.Float64
		}
	}

	return stats, nil
}

// UpdateStudyStreak updates the user's study streak
func (s *ProgressService) UpdateStudyStreak(ctx context.Context, userID int) error {
	today := time.Now().Format("2006-01-02")

	// Get current statistics
	var lastStudyDate sql.NullString
	var currentStreak int

	query := `
		SELECT last_study_date, study_streak_days
		FROM user_statistics
		WHERE user_id = $1
	`

	err := s.db.QueryRowContext(ctx, query, userID).Scan(&lastStudyDate, &currentStreak)
	if err == sql.ErrNoRows {
		// Create initial statistics record
		createQuery := `
			INSERT INTO user_statistics (user_id, last_study_date, study_streak_days)
			VALUES ($1, $2, 1)
			ON CONFLICT (user_id) DO UPDATE
			SET last_study_date = $2, study_streak_days = 1, updated_at = CURRENT_TIMESTAMP
		`
		_, err = s.db.ExecContext(ctx, createQuery, userID, today)
		return err
	} else if err != nil {
		return fmt.Errorf("error getting study streak: %w", err)
	}

	// Check if already studied today
	if lastStudyDate.Valid && lastStudyDate.String == today {
		return nil // Already counted for today
	}

	// Calculate new streak
	newStreak := 1
	if lastStudyDate.Valid {
		lastDate, err := time.Parse("2006-01-02", lastStudyDate.String)
		if err == nil {
			todayDate, _ := time.Parse("2006-01-02", today)
			daysDiff := int(todayDate.Sub(lastDate).Hours() / 24)

			if daysDiff == 1 {
				// Consecutive day
				newStreak = currentStreak + 1
			} else if daysDiff > 1 {
				// Streak broken
				newStreak = 1
			}
		}
	}

	// Update statistics
	updateQuery := `
		UPDATE user_statistics
		SET last_study_date = $1, study_streak_days = $2, updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $3
	`
	_, err = s.db.ExecContext(ctx, updateQuery, today, newStreak, userID)
	return err
}
