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

type vocabularyRepository struct {
	db *database.DB
}

func NewVocabularyRepository(db *database.DB) repository.VocabularyRepository {
	return &vocabularyRepository{db: db}
}

func (r *vocabularyRepository) GetAll(ctx context.Context, jlptLevel *int, limit, offset int) ([]models.Vocabulary, error) {
	query := `
		SELECT id, word, reading, meaning, part_of_speech, jlpt_level,
		       example_sentence, example_translation, audio_url, created_at, updated_at
		FROM vocabulary
		WHERE ($1::int IS NULL OR jlpt_level = $1)
		ORDER BY id
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, jlptLevel, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying vocabulary: %w", err)
	}
	defer rows.Close()

	var items []models.Vocabulary
	for rows.Next() {
		var v models.Vocabulary
		err := rows.Scan(
			&v.ID, &v.Word, &v.Reading, &v.Meaning, &v.PartOfSpeech, &v.JLPTLevel,
			&v.ExampleSentence, &v.ExampleTranslation, &v.AudioURL, &v.CreatedAt, &v.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning vocabulary: %w", err)
		}
		items = append(items, v)
	}

	return items, rows.Err()
}

func (r *vocabularyRepository) GetByID(ctx context.Context, id int) (*models.Vocabulary, error) {
	query := `
		SELECT id, word, reading, meaning, part_of_speech, jlpt_level,
		       example_sentence, example_translation, audio_url, created_at, updated_at
		FROM vocabulary
		WHERE id = $1
	`

	v := &models.Vocabulary{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&v.ID, &v.Word, &v.Reading, &v.Meaning, &v.PartOfSpeech, &v.JLPTLevel,
		&v.ExampleSentence, &v.ExampleTranslation, &v.AudioURL, &v.CreatedAt, &v.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, pkgErrors.NotFound("Vocabulary not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error getting vocabulary: %w", err)
	}

	return v, nil
}

func (r *vocabularyRepository) GetDueForReview(ctx context.Context, userID int, limit int) ([]models.VocabularyWithProgress, error) {
	query := `
		SELECT v.id, v.word, v.reading, v.meaning, v.part_of_speech, v.jlpt_level,
		       v.example_sentence, v.example_translation, v.audio_url, v.created_at, v.updated_at,
		       p.id, p.user_id, p.vocabulary_id, p.ease_factor, p.interval, p.repetitions,
		       p.next_review_date, p.last_reviewed_at, p.total_reviews, p.correct_reviews,
		       p.created_at, p.updated_at
		FROM vocabulary v
		INNER JOIN user_vocabulary_progress p ON v.id = p.vocabulary_id
		WHERE p.user_id = $1 AND p.next_review_date <= CURRENT_TIMESTAMP
		ORDER BY p.next_review_date
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("error querying due vocabulary: %w", err)
	}
	defer rows.Close()

	var items []models.VocabularyWithProgress
	for rows.Next() {
		var item models.VocabularyWithProgress
		item.Progress = &models.UserVocabularyProgress{}

		err := rows.Scan(
			&item.ID, &item.Word, &item.Reading, &item.Meaning, &item.PartOfSpeech, &item.JLPTLevel,
			&item.ExampleSentence, &item.ExampleTranslation, &item.AudioURL, &item.CreatedAt, &item.UpdatedAt,
			&item.Progress.ID, &item.Progress.UserID, &item.Progress.VocabularyID,
			&item.Progress.EaseFactor, &item.Progress.Interval, &item.Progress.Repetitions,
			&item.Progress.NextReviewDate, &item.Progress.LastReviewedAt, &item.Progress.TotalReviews,
			&item.Progress.CorrectReviews, &item.Progress.CreatedAt, &item.Progress.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning due vocabulary: %w", err)
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *vocabularyRepository) GetUserProgress(ctx context.Context, userID, vocabularyID int) (*models.UserVocabularyProgress, error) {
	query := `
		SELECT id, user_id, vocabulary_id, ease_factor, interval, repetitions,
		       next_review_date, last_reviewed_at, total_reviews, correct_reviews,
		       created_at, updated_at
		FROM user_vocabulary_progress
		WHERE user_id = $1 AND vocabulary_id = $2
	`

	p := &models.UserVocabularyProgress{}
	err := r.db.QueryRowContext(ctx, query, userID, vocabularyID).Scan(
		&p.ID, &p.UserID, &p.VocabularyID, &p.EaseFactor, &p.Interval, &p.Repetitions,
		&p.NextReviewDate, &p.LastReviewedAt, &p.TotalReviews, &p.CorrectReviews,
		&p.CreatedAt, &p.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, pkgErrors.NotFound("Progress not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error getting user progress: %w", err)
	}

	return p, nil
}

func (r *vocabularyRepository) CreateUserProgress(ctx context.Context, progress *models.UserVocabularyProgress) error {
	query := `
		INSERT INTO user_vocabulary_progress
		(user_id, vocabulary_id, ease_factor, interval, repetitions, next_review_date, total_reviews, correct_reviews)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		progress.UserID, progress.VocabularyID, progress.EaseFactor,
		progress.Interval, progress.Repetitions, progress.NextReviewDate,
		progress.TotalReviews, progress.CorrectReviews,
	).Scan(&progress.ID, &progress.CreatedAt, &progress.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error creating user progress: %w", err)
	}

	return nil
}

func (r *vocabularyRepository) UpdateUserProgress(ctx context.Context, progress *models.UserVocabularyProgress) error {
	query := `
		UPDATE user_vocabulary_progress
		SET ease_factor = $1, interval = $2, repetitions = $3, next_review_date = $4,
		    last_reviewed_at = $5, total_reviews = $6, correct_reviews = $7, updated_at = CURRENT_TIMESTAMP
		WHERE id = $8
	`

	result, err := r.db.ExecContext(ctx, query,
		progress.EaseFactor, progress.Interval, progress.Repetitions, progress.NextReviewDate,
		progress.LastReviewedAt, progress.TotalReviews, progress.CorrectReviews, progress.ID,
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

func (r *vocabularyRepository) GetUserVocabularyList(ctx context.Context, userID int, jlptLevel *int, limit, offset int) ([]models.VocabularyWithProgress, error) {
	query := `
		SELECT v.id, v.word, v.reading, v.meaning, v.part_of_speech, v.jlpt_level,
		       v.example_sentence, v.example_translation, v.audio_url, v.created_at, v.updated_at,
		       p.id, p.user_id, p.vocabulary_id, p.ease_factor, p.interval, p.repetitions,
		       p.next_review_date, p.last_reviewed_at, p.total_reviews, p.correct_reviews,
		       p.created_at, p.updated_at
		FROM vocabulary v
		LEFT JOIN user_vocabulary_progress p ON v.id = p.vocabulary_id AND p.user_id = $1
		WHERE ($2::int IS NULL OR v.jlpt_level = $2)
		ORDER BY v.id
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.QueryContext(ctx, query, userID, jlptLevel, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying user vocabulary list: %w", err)
	}
	defer rows.Close()

	var items []models.VocabularyWithProgress
	for rows.Next() {
		var item models.VocabularyWithProgress
		var progressID sql.NullInt64
		var progressUserID, progressVocabID sql.NullInt64
		var easeFactor sql.NullFloat64
		var interval, repetitions, totalReviews, correctReviews sql.NullInt64
		var nextReviewDate, lastReviewedAt, progressCreatedAt, progressUpdatedAt sql.NullTime

		err := rows.Scan(
			&item.ID, &item.Word, &item.Reading, &item.Meaning, &item.PartOfSpeech, &item.JLPTLevel,
			&item.ExampleSentence, &item.ExampleTranslation, &item.AudioURL, &item.CreatedAt, &item.UpdatedAt,
			&progressID, &progressUserID, &progressVocabID, &easeFactor, &interval, &repetitions,
			&nextReviewDate, &lastReviewedAt, &totalReviews, &correctReviews,
			&progressCreatedAt, &progressUpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning user vocabulary: %w", err)
		}

		// Only populate progress if it exists
		if progressID.Valid {
			item.Progress = &models.UserVocabularyProgress{
				ID:             int(progressID.Int64),
				UserID:         int(progressUserID.Int64),
				VocabularyID:   int(progressVocabID.Int64),
				EaseFactor:     easeFactor.Float64,
				Interval:       int(interval.Int64),
				Repetitions:    int(repetitions.Int64),
				NextReviewDate: nextReviewDate.Time,
				TotalReviews:   int(totalReviews.Int64),
				CorrectReviews: int(correctReviews.Int64),
				CreatedAt:      progressCreatedAt.Time,
				UpdatedAt:      progressUpdatedAt.Time,
			}
			if lastReviewedAt.Valid {
				item.Progress.LastReviewedAt = &lastReviewedAt.Time
			}
		}

		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *vocabularyRepository) Count(ctx context.Context, jlptLevel *int) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM vocabulary
		WHERE ($1::int IS NULL OR jlpt_level = $1)
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, jlptLevel).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting vocabulary: %w", err)
	}

	return count, nil
}
