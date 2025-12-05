-- Drop indexes for user grammar progress
DROP INDEX IF EXISTS idx_ugp_user_completed;
DROP INDEX IF EXISTS idx_ugp_grammar_lesson_id;
DROP INDEX IF EXISTS idx_ugp_user_id;

-- Drop user grammar progress table
DROP TABLE IF EXISTS user_grammar_progress;

-- Drop indexes for grammar examples
DROP INDEX IF EXISTS idx_grammar_examples_order;
DROP INDEX IF EXISTS idx_grammar_examples_lesson_id;

-- Drop grammar examples table
DROP TABLE IF EXISTS grammar_examples;

-- Drop indexes for grammar lessons
DROP INDEX IF EXISTS idx_grammar_lesson_order;
DROP INDEX IF EXISTS idx_grammar_jlpt_level;

-- Drop grammar lessons table
DROP TABLE IF EXISTS grammar_lessons;

-- Remove migration record
DELETE FROM schema_migrations WHERE version = '003_create_grammar_tables';
