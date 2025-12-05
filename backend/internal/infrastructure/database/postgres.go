package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/joaosantos/jlpt5/internal/config"
	"github.com/joaosantos/jlpt5/internal/utils"

	_ "github.com/lib/pq"
)

// DB wraps the database connection
type DB struct {
	*sql.DB
	logger *utils.Logger
}

// NewPostgresConnection creates a new PostgreSQL database connection
func NewPostgresConnection(cfg *config.DatabaseConfig, logger *utils.Logger) (*DB, error) {
	logger.Info("Connecting to PostgreSQL database", utils.WithContext(
		"host", cfg.Host,
		"port", cfg.Port,
		"database", cfg.DBName,
	))

	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	// Test the connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	logger.Info("Successfully connected to PostgreSQL database")

	return &DB{
		DB:     db,
		logger: logger,
	}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	db.logger.Info("Closing database connection")
	return db.DB.Close()
}

// HealthCheck checks if the database connection is healthy
func (db *DB) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

// RunMigrations runs database migrations
func (db *DB) RunMigrations() error {
	db.logger.Info("Running database migrations")

	// Create migrations table if it doesn't exist
	createMigrationsTable := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
	`

	if _, err := db.Exec(createMigrationsTable); err != nil {
		return fmt.Errorf("error creating migrations table: %w", err)
	}

	// Get list of already applied migrations
	appliedMigrations, err := db.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("error getting applied migrations: %w", err)
	}

	// Find migrations directory
	migrationsPath := filepath.Join("internal", "infrastructure", "postgres", "migrations")

	// Read migration files from filesystem
	entries, err := os.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("error reading migrations directory: %w", err)
	}

	// Filter and sort .up.sql files
	var migrationFileNames []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".up.sql") {
			migrationFileNames = append(migrationFileNames, entry.Name())
		}
	}
	sort.Strings(migrationFileNames)

	// Execute pending migrations
	for _, filename := range migrationFileNames {
		// Extract version from filename (e.g., "001_create_users_table.up.sql" -> "001_create_users_table")
		version := strings.TrimSuffix(filename, ".up.sql")

		if appliedMigrations[version] {
			db.logger.Debug("Skipping already applied migration", utils.WithContext("version", version))
			continue
		}

		db.logger.Info("Applying migration", utils.WithContext("version", version))

		// Read migration file from filesystem
		content, err := os.ReadFile(filepath.Join(migrationsPath, filename))
		if err != nil {
			return fmt.Errorf("error reading migration file %s: %w", filename, err)
		}

		// Execute migration in a transaction
		if err := db.executeMigration(version, string(content)); err != nil {
			return fmt.Errorf("error executing migration %s: %w", version, err)
		}

		db.logger.Info("Successfully applied migration", utils.WithContext("version", version))
	}

	db.logger.Info("All migrations completed successfully")
	return nil
}

// getAppliedMigrations returns a map of applied migration versions
func (db *DB) getAppliedMigrations() (map[string]bool, error) {
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, rows.Err()
}

// executeMigration executes a single migration in a transaction
func (db *DB) executeMigration(version, sql string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute the migration SQL
	if _, err := tx.Exec(sql); err != nil {
		return err
	}

	// The migration file itself inserts into schema_migrations,
	// so we don't need to do it here

	return tx.Commit()
}

// BeginTx starts a new transaction
func (db *DB) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return db.DB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
}

// WithTransaction executes a function within a database transaction
func (db *DB) WithTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			db.logger.Error("error rolling back transaction", utils.WithContext(
				"error", rbErr.Error(),
				"original_error", err.Error(),
			))
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}
