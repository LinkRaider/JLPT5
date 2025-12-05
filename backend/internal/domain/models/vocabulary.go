package models

import "time"

// Vocabulary represents a vocabulary item
type Vocabulary struct {
	ID                 int       `json:"id"`
	Word               string    `json:"word"`
	Reading            string    `json:"reading"`
	Meaning            string    `json:"meaning"`
	PartOfSpeech       *string   `json:"part_of_speech,omitempty"`
	JLPTLevel          int       `json:"jlpt_level"`
	ExampleSentence    *string   `json:"example_sentence,omitempty"`
	ExampleTranslation *string   `json:"example_translation,omitempty"`
	AudioURL           *string   `json:"audio_url,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// UserVocabularyProgress represents a user's progress with a vocabulary item
// This implements the SM-2 spaced repetition algorithm
type UserVocabularyProgress struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	VocabularyID   int       `json:"vocabulary_id"`
	EaseFactor     float64   `json:"ease_factor"`     // SM-2: Typically starts at 2.5
	Interval       int       `json:"interval"`        // Days until next review
	Repetitions    int       `json:"repetitions"`     // Number of successful consecutive reviews
	NextReviewDate time.Time `json:"next_review_date"`
	LastReviewedAt *time.Time `json:"last_reviewed_at,omitempty"`
	TotalReviews   int       `json:"total_reviews"`
	CorrectReviews int       `json:"correct_reviews"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// VocabularyWithProgress combines vocabulary with user progress
type VocabularyWithProgress struct {
	Vocabulary
	Progress *UserVocabularyProgress `json:"progress,omitempty"`
}
