package dto

// QuizResponse represents a quiz in API responses
type QuizResponse struct {
	ID           int     `json:"id"`
	Title        string  `json:"title"`
	Description  *string `json:"description,omitempty"`
	JLPTLevel    int     `json:"jlpt_level"`
	QuizType     *string `json:"quiz_type,omitempty"`
	PassingScore int     `json:"passing_score"`
}

// QuizListResponse represents a paginated list of quizzes
type QuizListResponse struct {
	Items      []QuizResponse `json:"items"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// QuizQuestionResponse represents a quiz question
type QuizQuestionResponse struct {
	ID           int      `json:"id"`
	QuestionText string   `json:"question_text"`
	QuestionType string   `json:"question_type"`
	Options      []string `json:"options"`
	Points       int      `json:"points"`
	// Note: CorrectAnswer and Explanation are excluded from this response
	// They are only included in results after submission
}

// QuizQuestionDetailResponse represents a quiz question with answer details (for results)
type QuizQuestionDetailResponse struct {
	ID            int      `json:"id"`
	QuestionText  string   `json:"question_text"`
	QuestionType  string   `json:"question_type"`
	Options       []string `json:"options"`
	CorrectAnswer string   `json:"correct_answer"`
	Explanation   *string  `json:"explanation,omitempty"`
	Points        int      `json:"points"`
}

// StartQuizResponse represents the response when starting a quiz
type StartQuizResponse struct {
	SessionID int                    `json:"session_id"`
	Quiz      QuizResponse           `json:"quiz"`
	Questions []QuizQuestionResponse `json:"questions"`
	StartedAt string                 `json:"started_at"`
}

// SubmitQuizRequest represents a request to submit quiz answers
type SubmitQuizRequest struct {
	Answers map[int]string `json:"answers"` // questionID -> userAnswer
}

// QuizAnswerResponse represents a submitted answer with correctness
type QuizAnswerResponse struct {
	QuestionID int    `json:"question_id"`
	UserAnswer string `json:"user_answer"`
	IsCorrect  bool   `json:"is_correct"`
}

// QuizResultResponse represents the result of a completed quiz
type QuizResultResponse struct {
	SessionID      int                          `json:"session_id"`
	Quiz           QuizResponse                 `json:"quiz"`
	Score          int                          `json:"score"`          // Points earned
	TotalPoints    int                          `json:"total_points"`   // Total points possible
	Percentage     float64                      `json:"percentage"`     // Score percentage
	TotalQuestions int                          `json:"total_questions"` // Number of questions
	Passed         bool                         `json:"passed"`
	StartedAt      string                       `json:"started_at"`
	CompletedAt    string                       `json:"completed_at"`
	Questions      []QuizQuestionDetailResponse `json:"questions"`
	Answers        []QuizAnswerResponse         `json:"answers"`
}

// QuizSessionResponse represents a quiz session summary
type QuizSessionResponse struct {
	ID          int      `json:"id"`
	QuizID      int      `json:"quiz_id"`
	StartedAt   string   `json:"started_at"`
	CompletedAt *string  `json:"completed_at,omitempty"`
	Score       *int     `json:"score,omitempty"`       // Points earned
	Percentage  *float64 `json:"percentage,omitempty"`  // Score percentage
	Passed      *bool    `json:"passed,omitempty"`
}

// QuizHistoryResponse represents a user's quiz history
type QuizHistoryResponse struct {
	Sessions []QuizSessionResponse `json:"sessions"`
	Page     int                   `json:"page"`
	PageSize int                   `json:"page_size"`
}
