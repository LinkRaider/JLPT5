package dto

// VocabularyResponse represents a vocabulary item in API responses
type VocabularyResponse struct {
	ID                 int                    `json:"id"`
	Word               string                 `json:"word"`
	Reading            string                 `json:"reading"`
	Meaning            string                 `json:"meaning"`
	PartOfSpeech       *string                `json:"part_of_speech,omitempty"`
	JLPTLevel          int                    `json:"jlpt_level"`
	ExampleSentence    *string                `json:"example_sentence,omitempty"`
	ExampleTranslation *string                `json:"example_translation,omitempty"`
	AudioURL           *string                `json:"audio_url,omitempty"`
	Progress           *ProgressResponse      `json:"progress,omitempty"`
}

// ProgressResponse represents user progress for a vocabulary item
type ProgressResponse struct {
	ID             int     `json:"id"`
	EaseFactor     float64 `json:"ease_factor"`
	Interval       int     `json:"interval_days"`
	Repetitions    int     `json:"repetitions"`
	NextReviewDate string  `json:"next_review_date"`
	LastReviewedAt *string `json:"last_reviewed_at,omitempty"`
	TotalReviews   int     `json:"total_reviews"`
	CorrectReviews int     `json:"correct_reviews"`
	SuccessRate    float64 `json:"success_rate"`
	IsDue          bool    `json:"is_due"`
}

// VocabularyListResponse represents a paginated list of vocabulary
type VocabularyListResponse struct {
	Items      []VocabularyResponse `json:"items"`
	Total      int                  `json:"total"`
	Page       int                  `json:"page"`
	PageSize   int                  `json:"page_size"`
	TotalPages int                  `json:"total_pages"`
}

// ReviewRequest represents a vocabulary review submission
type ReviewRequest struct {
	IsCorrect bool `json:"is_correct"`
}

// ReviewResponse represents the result of a review submission
type ReviewResponse struct {
	Success        bool              `json:"success"`
	Progress       *ProgressResponse `json:"progress"`
	NextReviewDate string            `json:"next_review_date"`
	Message        string            `json:"message"`
}
