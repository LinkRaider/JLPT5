-- Drop indexes for daily study log
DROP INDEX IF EXISTS idx_dsl_user_date;
DROP INDEX IF EXISTS idx_dsl_study_date;
DROP INDEX IF EXISTS idx_dsl_user_id;

-- Drop daily study log table
DROP TABLE IF EXISTS daily_study_log;

-- Drop indexes for user statistics
DROP INDEX IF EXISTS idx_user_statistics_study_streak;
DROP INDEX IF EXISTS idx_user_statistics_user_id;

-- Drop user statistics table
DROP TABLE IF EXISTS user_statistics;

-- Remove migration record
DELETE FROM schema_migrations WHERE version = '005_create_progress_tables';
