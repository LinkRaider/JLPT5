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

type grammarRepository struct {
	db *database.DB
}

func NewGrammarRepository(db *database.DB) repository.GrammarRepository {
	return &grammarRepository{db: db}
}

func (r *grammarRepository) GetAllLessons(ctx context.Context, jlptLevel *int, limit, offset int) ([]models.GrammarLesson, error) {
	query := `
		SELECT id, title, grammar_point, explanation, usage_notes, jlpt_level, lesson_order, created_at, updated_at
		FROM grammar_lessons
		WHERE ($1::int IS NULL OR jlpt_level = $1)
		ORDER BY lesson_order, id
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, jlptLevel, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying grammar lessons: %w", err)
	}
	defer rows.Close()

	var lessons []models.GrammarLesson
	for rows.Next() {
		var l models.GrammarLesson
		err := rows.Scan(&l.ID, &l.Title, &l.GrammarPoint, &l.Explanation, &l.UsageNotes,
			&l.JLPTLevel, &l.LessonOrder, &l.CreatedAt, &l.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning grammar lesson: %w", err)
		}
		lessons = append(lessons, l)
	}

	return lessons, rows.Err()
}

func (r *grammarRepository) GetLessonByID(ctx context.Context, lessonID int) (*models.GrammarLessonWithExamples, error) {
	// Get lesson
	lessonQuery := `
		SELECT id, title, grammar_point, explanation, usage_notes, jlpt_level, lesson_order, created_at, updated_at
		FROM grammar_lessons
		WHERE id = $1
	`

	lesson := &models.GrammarLessonWithExamples{}
	err := r.db.QueryRowContext(ctx, lessonQuery, lessonID).Scan(
		&lesson.ID, &lesson.Title, &lesson.GrammarPoint, &lesson.Explanation, &lesson.UsageNotes,
		&lesson.JLPTLevel, &lesson.LessonOrder, &lesson.CreatedAt, &lesson.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, pkgErrors.NotFound("Grammar lesson not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error getting grammar lesson: %w", err)
	}

	// Get examples
	examplesQuery := `
		SELECT id, grammar_lesson_id, japanese_sentence, english_translation, notes, example_order, created_at
		FROM grammar_examples
		WHERE grammar_lesson_id = $1
		ORDER BY example_order, id
	`

	rows, err := r.db.QueryContext(ctx, examplesQuery, lessonID)
	if err != nil {
		return nil, fmt.Errorf("error querying grammar examples: %w", err)
	}
	defer rows.Close()

	var examples []models.GrammarExample
	for rows.Next() {
		var e models.GrammarExample
		err := rows.Scan(&e.ID, &e.GrammarLessonID, &e.JapaneseSentence, &e.EnglishTranslation,
			&e.Notes, &e.ExampleOrder, &e.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning grammar example: %w", err)
		}
		examples = append(examples, e)
	}

	lesson.Examples = examples
	return lesson, rows.Err()
}

func (r *grammarRepository) GetUserProgress(ctx context.Context, userID, lessonID int) (*models.UserGrammarProgress, error) {
	query := `
		SELECT id, user_id, grammar_lesson_id, completed, completed_at, notes, created_at, updated_at
		FROM user_grammar_progress
		WHERE user_id = $1 AND grammar_lesson_id = $2
	`

	p := &models.UserGrammarProgress{}
	err := r.db.QueryRowContext(ctx, query, userID, lessonID).Scan(
		&p.ID, &p.UserID, &p.GrammarLessonID, &p.Completed, &p.CompletedAt,
		&p.Notes, &p.CreatedAt, &p.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, pkgErrors.NotFound("Progress not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error getting user progress: %w", err)
	}

	return p, nil
}

func (r *grammarRepository) CreateUserProgress(ctx context.Context, progress *models.UserGrammarProgress) error {
	query := `
		INSERT INTO user_grammar_progress (user_id, grammar_lesson_id, completed, completed_at, notes)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		progress.UserID, progress.GrammarLessonID, progress.Completed,
		progress.CompletedAt, progress.Notes,
	).Scan(&progress.ID, &progress.CreatedAt, &progress.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error creating user progress: %w", err)
	}

	return nil
}

func (r *grammarRepository) UpdateUserProgress(ctx context.Context, progress *models.UserGrammarProgress) error {
	query := `
		UPDATE user_grammar_progress
		SET completed = $1, completed_at = $2, notes = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
	`

	result, err := r.db.ExecContext(ctx, query,
		progress.Completed, progress.CompletedAt, progress.Notes, progress.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating user progress: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rows == 0 {
		return pkgErrors.NotFound("Progress not found")
	}

	return nil
}

func (r *grammarRepository) GetUserLessonsList(ctx context.Context, userID int, jlptLevel *int, limit, offset int) ([]models.GrammarLessonWithExamples, error) {
	query := `
		SELECT gl.id, gl.title, gl.grammar_point, gl.explanation, gl.usage_notes,
		       gl.jlpt_level, gl.lesson_order, gl.created_at, gl.updated_at,
		       p.id, p.user_id, p.grammar_lesson_id, p.completed, p.completed_at,
		       p.notes, p.created_at, p.updated_at
		FROM grammar_lessons gl
		LEFT JOIN user_grammar_progress p ON gl.id = p.grammar_lesson_id AND p.user_id = $1
		WHERE ($2::int IS NULL OR gl.jlpt_level = $2)
		ORDER BY gl.lesson_order, gl.id
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.QueryContext(ctx, query, userID, jlptLevel, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying user lessons list: %w", err)
	}
	defer rows.Close()

	var lessons []models.GrammarLessonWithExamples
	for rows.Next() {
		var lesson models.GrammarLessonWithExamples
		var progressID sql.NullInt64
		var progressUserID, progressLessonID sql.NullInt64
		var completed sql.NullBool
		var completedAt sql.NullTime
		var notes sql.NullString
		var progressCreatedAt, progressUpdatedAt sql.NullTime

		err := rows.Scan(
			&lesson.ID, &lesson.Title, &lesson.GrammarPoint, &lesson.Explanation, &lesson.UsageNotes,
			&lesson.JLPTLevel, &lesson.LessonOrder, &lesson.CreatedAt, &lesson.UpdatedAt,
			&progressID, &progressUserID, &progressLessonID, &completed, &completedAt,
			&notes, &progressCreatedAt, &progressUpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning user lesson: %w", err)
		}

		// Populate progress if it exists
		if progressID.Valid {
			lesson.Progress = &models.UserGrammarProgress{
				ID:              int(progressID.Int64),
				UserID:          int(progressUserID.Int64),
				GrammarLessonID: int(progressLessonID.Int64),
				Completed:       completed.Bool,
				CreatedAt:       progressCreatedAt.Time,
				UpdatedAt:       progressUpdatedAt.Time,
			}
			if completedAt.Valid {
				lesson.Progress.CompletedAt = &completedAt.Time
			}
			if notes.Valid {
				lesson.Progress.Notes = &notes.String
			}
		}

		// Get examples for this lesson
		examplesQuery := `
			SELECT id, grammar_lesson_id, japanese_sentence, english_translation, notes, example_order, created_at
			FROM grammar_examples
			WHERE grammar_lesson_id = $1
			ORDER BY example_order, id
		`

		exampleRows, err := r.db.QueryContext(ctx, examplesQuery, lesson.ID)
		if err != nil {
			return nil, fmt.Errorf("error querying grammar examples: %w", err)
		}

		var examples []models.GrammarExample
		for exampleRows.Next() {
			var e models.GrammarExample
			err := exampleRows.Scan(&e.ID, &e.GrammarLessonID, &e.JapaneseSentence, &e.EnglishTranslation,
				&e.Notes, &e.ExampleOrder, &e.CreatedAt)
			if err != nil {
				exampleRows.Close()
				return nil, fmt.Errorf("error scanning grammar example: %w", err)
			}
			examples = append(examples, e)
		}
		exampleRows.Close()

		lesson.Examples = examples
		lessons = append(lessons, lesson)
	}

	return lessons, rows.Err()
}

func (r *grammarRepository) CountLessons(ctx context.Context, jlptLevel *int) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM grammar_lessons
		WHERE ($1::int IS NULL OR jlpt_level = $1)
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, jlptLevel).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting grammar lessons: %w", err)
	}

	return count, nil
}
