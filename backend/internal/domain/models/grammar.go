package models

import "time"

// GrammarLesson represents a grammar lesson
type GrammarLesson struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	GrammarPoint string    `json:"grammar_point"`
	Explanation  string    `json:"explanation"`
	UsageNotes   *string   `json:"usage_notes,omitempty"`
	JLPTLevel    int       `json:"jlpt_level"`
	LessonOrder  *int      `json:"lesson_order,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// GrammarExample represents an example for a grammar lesson
type GrammarExample struct {
	ID                 int       `json:"id"`
	GrammarLessonID    int       `json:"grammar_lesson_id"`
	JapaneseSentence   string    `json:"japanese_sentence"`
	EnglishTranslation string    `json:"english_translation"`
	Notes              *string   `json:"notes,omitempty"`
	ExampleOrder       *int      `json:"example_order,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
}

// UserGrammarProgress represents a user's progress with a grammar lesson
type UserGrammarProgress struct {
	ID              int        `json:"id"`
	UserID          int        `json:"user_id"`
	GrammarLessonID int        `json:"grammar_lesson_id"`
	Completed       bool       `json:"completed"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
	Notes           *string    `json:"notes,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// GrammarLessonWithExamples combines a lesson with its examples
type GrammarLessonWithExamples struct {
	GrammarLesson
	Examples []GrammarExample        `json:"examples"`
	Progress *UserGrammarProgress    `json:"progress,omitempty"`
}
