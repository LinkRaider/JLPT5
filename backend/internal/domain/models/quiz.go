package models

import "time"

// QuestionType represents the type of quiz question
type QuestionType string

const (
	QuestionTypeMultipleChoice QuestionType = "multiple_choice"
	QuestionTypeFillInBlank    QuestionType = "fill_in_blank"
)

// Quiz represents a quiz
type Quiz struct {
	ID              int       `json:"id"`
	Title           string    `json:"title"`
	Description     *string   `json:"description,omitempty"`
	QuizType        *string   `json:"quiz_type,omitempty"` // vocabulary, grammar, mixed
	JLPTLevel       int       `json:"jlpt_level"`
	TimeLimitMinutes *int     `json:"time_limit_minutes,omitempty"`
	PassingScore    int       `json:"passing_score"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// QuizQuestion represents a question in a quiz
type QuizQuestion struct {
	ID           int          `json:"id"`
	QuizID       int          `json:"quiz_id"`
	QuestionType QuestionType `json:"question_type"`
	QuestionText string       `json:"question_text"`
	CorrectAnswer string      `json:"-"` // Hidden from JSON responses
	OptionA      *string      `json:"option_a,omitempty"`
	OptionB      *string      `json:"option_b,omitempty"`
	OptionC      *string      `json:"option_c,omitempty"`
	OptionD      *string      `json:"option_d,omitempty"`
	Explanation  *string      `json:"explanation,omitempty"`
	Points       int          `json:"points"`
	QuestionOrder *int        `json:"question_order,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
}

// QuizSession represents a user's quiz attempt
type QuizSession struct {
	ID              int        `json:"id"`
	UserID          int        `json:"user_id"`
	QuizID          int        `json:"quiz_id"`
	StartedAt       time.Time  `json:"started_at"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
	Score           *int       `json:"score,omitempty"`
	TotalPoints     *int       `json:"total_points,omitempty"`
	Percentage      *float64   `json:"percentage,omitempty"`
	Passed          *bool      `json:"passed,omitempty"`
	TimeSpentSeconds *int      `json:"time_spent_seconds,omitempty"`
}

// QuizAnswer represents a user's answer to a quiz question
type QuizAnswer struct {
	ID             int       `json:"id"`
	QuizSessionID  int       `json:"quiz_session_id"`
	QuizQuestionID int       `json:"quiz_question_id"`
	UserAnswer     *string   `json:"user_answer,omitempty"`
	IsCorrect      *bool     `json:"is_correct,omitempty"`
	AnsweredAt     time.Time `json:"answered_at"`
}

// QuizWithQuestions combines a quiz with its questions
type QuizWithQuestions struct {
	Quiz
	Questions []QuizQuestion `json:"questions"`
}

// QuizResult represents the results of a completed quiz session
type QuizResult struct {
	Session  QuizSession               `json:"session"`
	Quiz     Quiz                      `json:"quiz"`
	Answers  []QuizAnswerWithQuestion  `json:"answers"`
}

// QuizAnswerWithQuestion combines an answer with the question details
type QuizAnswerWithQuestion struct {
	QuizAnswer
	Question QuizQuestion `json:"question"`
	CorrectAnswer string  `json:"correct_answer"`
}
