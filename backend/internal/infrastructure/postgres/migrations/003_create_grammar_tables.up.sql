-- Create grammar lessons table
CREATE TABLE IF NOT EXISTS grammar_lessons (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    grammar_point VARCHAR(100) NOT NULL,
    explanation TEXT NOT NULL,
    usage_notes TEXT,
    jlpt_level INTEGER DEFAULT 5,
    lesson_order INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for grammar lessons
CREATE INDEX idx_grammar_jlpt_level ON grammar_lessons(jlpt_level);
CREATE INDEX idx_grammar_lesson_order ON grammar_lessons(lesson_order);

-- Create grammar examples table
CREATE TABLE IF NOT EXISTS grammar_examples (
    id SERIAL PRIMARY KEY,
    grammar_lesson_id INTEGER NOT NULL REFERENCES grammar_lessons(id) ON DELETE CASCADE,
    japanese_sentence TEXT NOT NULL,
    english_translation TEXT NOT NULL,
    notes TEXT,
    example_order INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for grammar examples
CREATE INDEX idx_grammar_examples_lesson_id ON grammar_examples(grammar_lesson_id);
CREATE INDEX idx_grammar_examples_order ON grammar_examples(grammar_lesson_id, example_order);

-- Create user grammar progress table
CREATE TABLE IF NOT EXISTS user_grammar_progress (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    grammar_lesson_id INTEGER NOT NULL REFERENCES grammar_lessons(id) ON DELETE CASCADE,
    completed BOOLEAN DEFAULT FALSE,
    completed_at TIMESTAMP WITH TIME ZONE,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_user_grammar UNIQUE(user_id, grammar_lesson_id)
);

-- Create indexes for user grammar progress
CREATE INDEX idx_ugp_user_id ON user_grammar_progress(user_id);
CREATE INDEX idx_ugp_grammar_lesson_id ON user_grammar_progress(grammar_lesson_id);
CREATE INDEX idx_ugp_user_completed ON user_grammar_progress(user_id, completed);

-- Insert migration record
INSERT INTO schema_migrations (version) VALUES ('003_create_grammar_tables')
ON CONFLICT (version) DO NOTHING;
