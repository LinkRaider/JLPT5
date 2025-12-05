-- Drop indexes for quiz answers
DROP INDEX IF EXISTS idx_quiz_answers_question_id;
DROP INDEX IF EXISTS idx_quiz_answers_session_id;

-- Drop quiz answers table
DROP TABLE IF EXISTS quiz_answers;

-- Drop indexes for quiz sessions
DROP INDEX IF EXISTS idx_quiz_sessions_user_completed;
DROP INDEX IF EXISTS idx_quiz_sessions_quiz_id;
DROP INDEX IF EXISTS idx_quiz_sessions_user_id;

-- Drop quiz sessions table
DROP TABLE IF EXISTS quiz_sessions;

-- Drop indexes for quiz questions
DROP INDEX IF EXISTS idx_quiz_questions_order;
DROP INDEX IF EXISTS idx_quiz_questions_quiz_id;

-- Drop quiz questions table
DROP TABLE IF EXISTS quiz_questions;

-- Drop indexes for quizzes
DROP INDEX IF EXISTS idx_quizzes_type;
DROP INDEX IF EXISTS idx_quizzes_jlpt_level;

-- Drop quizzes table
DROP TABLE IF EXISTS quizzes;

-- Drop question type enum
DROP TYPE IF EXISTS question_type;

-- Remove migration record
DELETE FROM schema_migrations WHERE version = '004_create_quiz_tables';
