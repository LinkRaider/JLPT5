package dto

// GrammarLessonResponse represents a grammar lesson in API responses
type GrammarLessonResponse struct {
	ID           int                      `json:"id"`
	Title        string                   `json:"title"`
	GrammarPoint string                   `json:"grammar_point"`
	Explanation  string                   `json:"explanation"`
	UsageNotes   *string                  `json:"usage_notes,omitempty"`
	JLPTLevel    int                      `json:"jlpt_level"`
	LessonOrder  *int                     `json:"lesson_order,omitempty"`
	Examples     []GrammarExampleResponse `json:"examples"`
	Progress     *GrammarProgressResponse `json:"progress,omitempty"`
}

// GrammarExampleResponse represents a grammar example
type GrammarExampleResponse struct {
	ID                 int     `json:"id"`
	JapaneseSentence   string  `json:"japanese_sentence"`
	EnglishTranslation string  `json:"english_translation"`
	Notes              *string `json:"notes,omitempty"`
}

// GrammarProgressResponse represents user progress for a grammar lesson
type GrammarProgressResponse struct {
	ID          int     `json:"id"`
	Completed   bool    `json:"completed"`
	CompletedAt *string `json:"completed_at,omitempty"`
	Notes       *string `json:"notes,omitempty"`
}

// GrammarListResponse represents a paginated list of grammar lessons
type GrammarListResponse struct {
	Items      []GrammarLessonResponse `json:"items"`
	Total      int                     `json:"total"`
	Page       int                     `json:"page"`
	PageSize   int                     `json:"page_size"`
	TotalPages int                     `json:"total_pages"`
}

// MarkCompletedRequest represents a request to mark a lesson as completed
type MarkCompletedRequest struct {
	Notes *string `json:"notes,omitempty"`
}
