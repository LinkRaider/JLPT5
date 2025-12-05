package dto

// ProgressStatsResponse represents overall user progress statistics
type ProgressStatsResponse struct {
	UserID                int     `json:"user_id"`
	StudyStreakDays       int     `json:"study_streak_days"`
	LastStudyDate         *string `json:"last_study_date,omitempty"`
	TotalStudyTimeMinutes int     `json:"total_study_time_minutes"`
	VocabularyLearned     int     `json:"vocabulary_learned"`
	VocabularyMastered    int     `json:"vocabulary_mastered"`
	VocabularyDue         int     `json:"vocabulary_due"`
	GrammarCompleted      int     `json:"grammar_completed"`
	GrammarTotal          int     `json:"grammar_total"`
	QuizzesTaken          int     `json:"quizzes_taken"`
	QuizzesPassed         int     `json:"quizzes_passed"`
	AverageQuizScore      float64 `json:"average_quiz_score"`
}
