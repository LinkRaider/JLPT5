-- Create user statistics table
CREATE TABLE IF NOT EXISTS user_statistics (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE UNIQUE,
    study_streak_days INTEGER DEFAULT 0,
    last_study_date DATE,
    total_study_time_minutes INTEGER DEFAULT 0,
    vocabulary_learned INTEGER DEFAULT 0,
    grammar_completed INTEGER DEFAULT 0,
    quizzes_taken INTEGER DEFAULT 0,
    quizzes_passed INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for user statistics
CREATE INDEX idx_user_statistics_user_id ON user_statistics(user_id);
CREATE INDEX idx_user_statistics_study_streak ON user_statistics(study_streak_days DESC);

-- Create daily study log table
CREATE TABLE IF NOT EXISTS daily_study_log (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    study_date DATE NOT NULL,
    vocabulary_reviewed INTEGER DEFAULT 0,
    grammar_studied INTEGER DEFAULT 0,
    quizzes_completed INTEGER DEFAULT 0,
    study_time_minutes INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_user_study_date UNIQUE(user_id, study_date)
);

-- Create indexes for daily study log
CREATE INDEX idx_dsl_user_id ON daily_study_log(user_id);
CREATE INDEX idx_dsl_study_date ON daily_study_log(study_date DESC);
CREATE INDEX idx_dsl_user_date ON daily_study_log(user_id, study_date DESC);

-- Insert migration record
INSERT INTO schema_migrations (version) VALUES ('005_create_progress_tables')
ON CONFLICT (version) DO NOTHING;
