-- Drop indexes for user vocabulary progress
DROP INDEX IF EXISTS idx_uvp_next_review;
DROP INDEX IF EXISTS idx_uvp_vocabulary_id;
DROP INDEX IF EXISTS idx_uvp_user_id;

-- Drop user vocabulary progress table
DROP TABLE IF EXISTS user_vocabulary_progress;

-- Drop indexes for vocabulary
DROP INDEX IF EXISTS idx_vocabulary_part_of_speech;
DROP INDEX IF EXISTS idx_vocabulary_word;
DROP INDEX IF EXISTS idx_vocabulary_jlpt_level;

-- Drop vocabulary table
DROP TABLE IF EXISTS vocabulary;

-- Remove migration record
DELETE FROM schema_migrations WHERE version = '002_create_vocabulary_tables';
