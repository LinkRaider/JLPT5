-- Create question type enum
CREATE TYPE question_type AS ENUM ('multiple_choice', 'fill_in_blank');

-- Create quizzes table
CREATE TABLE IF NOT EXISTS quizzes (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    quiz_type VARCHAR(50),  -- vocabulary, grammar, mixed
    jlpt_level INTEGER DEFAULT 5,
    time_limit_minutes INTEGER,
    passing_score INTEGER DEFAULT 70,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for quizzes
CREATE INDEX idx_quizzes_jlpt_level ON quizzes(jlpt_level);
CREATE INDEX idx_quizzes_type ON quizzes(quiz_type);

-- Create quiz questions table
CREATE TABLE IF NOT EXISTS quiz_questions (
    id SERIAL PRIMARY KEY,
    quiz_id INTEGER NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
    question_type question_type NOT NULL,
    question_text TEXT NOT NULL,
    correct_answer TEXT NOT NULL,
    option_a TEXT,
    option_b TEXT,
    option_c TEXT,
    option_d TEXT,
    explanation TEXT,
    points INTEGER DEFAULT 1,
    question_order INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for quiz questions
CREATE INDEX idx_quiz_questions_quiz_id ON quiz_questions(quiz_id);
CREATE INDEX idx_quiz_questions_order ON quiz_questions(quiz_id, question_order);

-- Create quiz sessions table
CREATE TABLE IF NOT EXISTS quiz_sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    quiz_id INTEGER NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
    started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP WITH TIME ZONE,
    score INTEGER,
    total_points INTEGER,
    percentage DECIMAL(5,2),
    passed BOOLEAN,
    time_spent_seconds INTEGER
);

-- Create indexes for quiz sessions
CREATE INDEX idx_quiz_sessions_user_id ON quiz_sessions(user_id);
CREATE INDEX idx_quiz_sessions_quiz_id ON quiz_sessions(quiz_id);
CREATE INDEX idx_quiz_sessions_user_completed ON quiz_sessions(user_id, completed_at);

-- Create quiz answers table
CREATE TABLE IF NOT EXISTS quiz_answers (
    id SERIAL PRIMARY KEY,
    quiz_session_id INTEGER NOT NULL REFERENCES quiz_sessions(id) ON DELETE CASCADE,
    quiz_question_id INTEGER NOT NULL REFERENCES quiz_questions(id) ON DELETE CASCADE,
    user_answer TEXT,
    is_correct BOOLEAN,
    answered_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for quiz answers
CREATE INDEX idx_quiz_answers_session_id ON quiz_answers(quiz_session_id);
CREATE INDEX idx_quiz_answers_question_id ON quiz_answers(quiz_question_id);

-- Insert migration record
INSERT INTO schema_migrations (version) VALUES ('004_create_quiz_tables')
ON CONFLICT (version) DO NOTHING;
