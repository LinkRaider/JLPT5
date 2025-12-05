package models

import "time"

// User represents a user account
type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // Never include in JSON responses
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	IsActive     bool      `json:"is_active"`
}

// UserStatistics represents overall user statistics
type UserStatistics struct {
	ID                   int       `json:"id"`
	UserID               int       `json:"user_id"`
	StudyStreakDays      int       `json:"study_streak_days"`
	LastStudyDate        *time.Time `json:"last_study_date,omitempty"`
	TotalStudyTimeMinutes int      `json:"total_study_time_minutes"`
	VocabularyLearned    int       `json:"vocabulary_learned"`
	GrammarCompleted     int       `json:"grammar_completed"`
	QuizzesTaken         int       `json:"quizzes_taken"`
	QuizzesPassed        int       `json:"quizzes_passed"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// DailyStudyLog represents daily study activity
type DailyStudyLog struct {
	ID                 int       `json:"id"`
	UserID             int       `json:"user_id"`
	StudyDate          time.Time `json:"study_date"`
	VocabularyReviewed int       `json:"vocabulary_reviewed"`
	GrammarStudied     int       `json:"grammar_studied"`
	QuizzesCompleted   int       `json:"quizzes_completed"`
	StudyTimeMinutes   int       `json:"study_time_minutes"`
	CreatedAt          time.Time `json:"created_at"`
}
