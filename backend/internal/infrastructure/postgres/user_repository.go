package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/joaosantos/jlpt5/internal/domain/models"
	"github.com/joaosantos/jlpt5/internal/domain/repository"
	"github.com/joaosantos/jlpt5/internal/infrastructure/database"
	pkgErrors "github.com/joaosantos/jlpt5/pkg/errors"
)

// userRepository implements the UserRepository interface
type userRepository struct {
	db *database.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *database.DB) repository.UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (email, username, password_hash, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query, user.Email, user.Username, user.PasswordHash, user.IsActive).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		// Check for unique constraint violation
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			return pkgErrors.Conflict("Email already exists")
		}
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_username_key\"" {
			return pkgErrors.Conflict("Username already exists")
		}
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	query := `
		SELECT id, email, username, password_hash, created_at, updated_at, last_login_at, is_active
		FROM users
		WHERE id = $1 AND is_active = true
	`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
		&user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt, &user.IsActive,
	)

	if err == sql.ErrNoRows {
		return nil, pkgErrors.NotFound("User not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error getting user by ID: %w", err)
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, username, password_hash, created_at, updated_at, last_login_at, is_active
		FROM users
		WHERE email = $1
	`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
		&user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt, &user.IsActive,
	)

	if err == sql.ErrNoRows {
		return nil, pkgErrors.NotFound("User not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error getting user by email: %w", err)
	}

	return user, nil
}

// GetByUsername retrieves a user by username
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, email, username, password_hash, created_at, updated_at, last_login_at, is_active
		FROM users
		WHERE username = $1
	`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
		&user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt, &user.IsActive,
	)

	if err == sql.ErrNoRows {
		return nil, pkgErrors.NotFound("User not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error getting user by username: %w", err)
	}

	return user, nil
}

// Update updates a user
func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET email = $1, username = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`

	result, err := r.db.ExecContext(ctx, query, user.Email, user.Username, user.ID)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rows == 0 {
		return pkgErrors.NotFound("User not found")
	}

	return nil
}

// UpdateLastLogin updates the user's last login timestamp
func (r *userRepository) UpdateLastLogin(ctx context.Context, userID int) error {
	query := `
		UPDATE users
		SET last_login_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("error updating last login: %w", err)
	}

	return nil
}

// GetStatistics retrieves user statistics
func (r *userRepository) GetStatistics(ctx context.Context, userID int) (*models.UserStatistics, error) {
	query := `
		SELECT id, user_id, study_streak_days, last_study_date, total_study_time_minutes,
		       vocabulary_learned, grammar_completed, quizzes_taken, quizzes_passed,
		       created_at, updated_at
		FROM user_statistics
		WHERE user_id = $1
	`

	stats := &models.UserStatistics{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&stats.ID, &stats.UserID, &stats.StudyStreakDays, &stats.LastStudyDate,
		&stats.TotalStudyTimeMinutes, &stats.VocabularyLearned, &stats.GrammarCompleted,
		&stats.QuizzesTaken, &stats.QuizzesPassed, &stats.CreatedAt, &stats.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, pkgErrors.NotFound("User statistics not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error getting user statistics: %w", err)
	}

	return stats, nil
}

// CreateStatistics creates initial statistics for a user
func (r *userRepository) CreateStatistics(ctx context.Context, userID int) error {
	query := `
		INSERT INTO user_statistics (user_id, last_study_date)
		VALUES ($1, $2)
	`

	_, err := r.db.ExecContext(ctx, query, userID, time.Now())
	if err != nil {
		return fmt.Errorf("error creating user statistics: %w", err)
	}

	return nil
}

// UpdateStatistics updates user statistics
func (r *userRepository) UpdateStatistics(ctx context.Context, stats *models.UserStatistics) error {
	query := `
		UPDATE user_statistics
		SET study_streak_days = $1, last_study_date = $2, total_study_time_minutes = $3,
		    vocabulary_learned = $4, grammar_completed = $5, quizzes_taken = $6,
		    quizzes_passed = $7, updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $8
	`

	result, err := r.db.ExecContext(ctx, query,
		stats.StudyStreakDays, stats.LastStudyDate, stats.TotalStudyTimeMinutes,
		stats.VocabularyLearned, stats.GrammarCompleted, stats.QuizzesTaken,
		stats.QuizzesPassed, stats.UserID,
	)
	if err != nil {
		return fmt.Errorf("error updating user statistics: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rows == 0 {
		return pkgErrors.NotFound("User statistics not found")
	}

	return nil
}
