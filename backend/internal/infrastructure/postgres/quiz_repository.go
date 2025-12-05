package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/joaosantos/jlpt5/internal/domain/models"
	"github.com/joaosantos/jlpt5/internal/domain/repository"
	"github.com/joaosantos/jlpt5/internal/infrastructure/database"
	pkgErrors "github.com/joaosantos/jlpt5/pkg/errors"
)

type quizRepository struct {
	db *database.DB
}

func NewQuizRepository(db *database.DB) repository.QuizRepository {
	return &quizRepository{db: db}
}

func (r *quizRepository) GetAllQuizzes(ctx context.Context, jlptLevel *int, limit, offset int) ([]models.Quiz, error) {
	query := `
		SELECT id, title, description, jlpt_level, quiz_type, passing_score, created_at, updated_at
		FROM quizzes
		WHERE ($1::int IS NULL OR jlpt_level = $1)
		ORDER BY id
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, jlptLevel, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying quizzes: %w", err)
	}
	defer rows.Close()

	var quizzes []models.Quiz
	for rows.Next() {
		var q models.Quiz
		err := rows.Scan(&q.ID, &q.Title, &q.Description, &q.JLPTLevel, &q.QuizType,
			&q.PassingScore, &q.CreatedAt, &q.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning quiz: %w", err)
		}
		quizzes = append(quizzes, q)
	}

	return quizzes, rows.Err()
}

func (r *quizRepository) GetQuizByID(ctx context.Context, quizID int) (*models.Quiz, error) {
	query := `
		SELECT id, title, description, jlpt_level, quiz_type, passing_score, created_at, updated_at
		FROM quizzes
		WHERE id = $1
	`

	quiz := &models.Quiz{}
	err := r.db.QueryRowContext(ctx, query, quizID).Scan(
		&quiz.ID, &quiz.Title, &quiz.Description, &quiz.JLPTLevel, &quiz.QuizType,
		&quiz.PassingScore, &quiz.CreatedAt, &quiz.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, pkgErrors.NotFound("Quiz not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error getting quiz: %w", err)
	}

	return quiz, nil
}

func (r *quizRepository) GetQuizQuestions(ctx context.Context, quizID int) ([]models.QuizQuestion, error) {
	query := `
		SELECT id, quiz_id, question_text, question_type, correct_answer, option_a,
		       option_b, option_c, option_d, explanation, points, question_order, created_at
		FROM quiz_questions
		WHERE quiz_id = $1
		ORDER BY question_order, id
	`

	rows, err := r.db.QueryContext(ctx, query, quizID)
	if err != nil {
		return nil, fmt.Errorf("error querying quiz questions: %w", err)
	}
	defer rows.Close()

	var questions []models.QuizQuestion
	for rows.Next() {
		var q models.QuizQuestion
		err := rows.Scan(&q.ID, &q.QuizID, &q.QuestionText, &q.QuestionType, &q.CorrectAnswer,
			&q.OptionA, &q.OptionB, &q.OptionC, &q.OptionD, &q.Explanation, &q.Points, &q.QuestionOrder, &q.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning quiz question: %w", err)
		}
		questions = append(questions, q)
	}

	return questions, rows.Err()
}

func (r *quizRepository) CreateQuizSession(ctx context.Context, session *models.QuizSession) error {
	query := `
		INSERT INTO quiz_sessions (user_id, quiz_id, started_at)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query,
		session.UserID, session.QuizID, session.StartedAt,
	).Scan(&session.ID)

	if err != nil {
		return fmt.Errorf("error creating quiz session: %w", err)
	}

	return nil
}

func (r *quizRepository) UpdateQuizSession(ctx context.Context, session *models.QuizSession) error {
	query := `
		UPDATE quiz_sessions
		SET completed_at = $1, score = $2, total_points = $3,
		    percentage = $4, passed = $5, time_spent_seconds = $6
		WHERE id = $7
	`

	result, err := r.db.ExecContext(ctx, query,
		session.CompletedAt, session.Score, session.TotalPoints,
		session.Percentage, session.Passed, session.TimeSpentSeconds, session.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating quiz session: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rows == 0 {
		return pkgErrors.NotFound("Quiz session not found")
	}

	return nil
}

func (r *quizRepository) GetQuizSession(ctx context.Context, sessionID int) (*models.QuizSession, error) {
	query := `
		SELECT id, user_id, quiz_id, started_at, completed_at, score,
		       total_points, percentage, passed, time_spent_seconds
		FROM quiz_sessions
		WHERE id = $1
	`

	session := &models.QuizSession{}
	err := r.db.QueryRowContext(ctx, query, sessionID).Scan(
		&session.ID, &session.UserID, &session.QuizID, &session.StartedAt,
		&session.CompletedAt, &session.Score, &session.TotalPoints,
		&session.Percentage, &session.Passed, &session.TimeSpentSeconds,
	)

	if err == sql.ErrNoRows {
		return nil, pkgErrors.NotFound("Quiz session not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error getting quiz session: %w", err)
	}

	return session, nil
}

func (r *quizRepository) GetUserQuizSessions(ctx context.Context, userID int, limit, offset int) ([]models.QuizSession, error) {
	query := `
		SELECT id, user_id, quiz_id, started_at, completed_at, score,
		       total_points, percentage, passed, time_spent_seconds
		FROM quiz_sessions
		WHERE user_id = $1
		ORDER BY started_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying user quiz sessions: %w", err)
	}
	defer rows.Close()

	var sessions []models.QuizSession
	for rows.Next() {
		var s models.QuizSession
		err := rows.Scan(&s.ID, &s.UserID, &s.QuizID, &s.StartedAt,
			&s.CompletedAt, &s.Score, &s.TotalPoints,
			&s.Percentage, &s.Passed, &s.TimeSpentSeconds)
		if err != nil {
			return nil, fmt.Errorf("error scanning quiz session: %w", err)
		}
		sessions = append(sessions, s)
	}

	return sessions, rows.Err()
}

func (r *quizRepository) CreateQuizAnswer(ctx context.Context, answer *models.QuizAnswer) error {
	query := `
		INSERT INTO quiz_answers (quiz_session_id, quiz_question_id, user_answer, is_correct, answered_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query,
		answer.QuizSessionID, answer.QuizQuestionID, answer.UserAnswer,
		answer.IsCorrect, answer.AnsweredAt,
	).Scan(&answer.ID)

	if err != nil {
		return fmt.Errorf("error creating quiz answer: %w", err)
	}

	return nil
}

func (r *quizRepository) GetSessionAnswers(ctx context.Context, sessionID int) ([]models.QuizAnswer, error) {
	query := `
		SELECT id, quiz_session_id, quiz_question_id, user_answer, is_correct, answered_at
		FROM quiz_answers
		WHERE quiz_session_id = $1
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("error querying session answers: %w", err)
	}
	defer rows.Close()

	var answers []models.QuizAnswer
	for rows.Next() {
		var a models.QuizAnswer
		err := rows.Scan(&a.ID, &a.QuizSessionID, &a.QuizQuestionID, &a.UserAnswer,
			&a.IsCorrect, &a.AnsweredAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning quiz answer: %w", err)
		}
		answers = append(answers, a)
	}

	return answers, rows.Err()
}

func (r *quizRepository) CountQuizzes(ctx context.Context, jlptLevel *int) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM quizzes
		WHERE ($1::int IS NULL OR jlpt_level = $1)
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, jlptLevel).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting quizzes: %w", err)
	}

	return count, nil
}
