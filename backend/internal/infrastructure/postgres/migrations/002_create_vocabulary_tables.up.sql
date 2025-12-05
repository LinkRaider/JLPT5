-- Create vocabulary table
CREATE TABLE IF NOT EXISTS vocabulary (
    id SERIAL PRIMARY KEY,
    word VARCHAR(100) NOT NULL,
    reading VARCHAR(100) NOT NULL,
    meaning TEXT NOT NULL,
    part_of_speech VARCHAR(50),
    jlpt_level INTEGER DEFAULT 5,
    example_sentence TEXT,
    example_translation TEXT,
    audio_url VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for vocabulary
CREATE INDEX idx_vocabulary_jlpt_level ON vocabulary(jlpt_level);
CREATE INDEX idx_vocabulary_word ON vocabulary(word);
CREATE INDEX idx_vocabulary_part_of_speech ON vocabulary(part_of_speech);

-- Create user vocabulary progress table (for spaced repetition)
CREATE TABLE IF NOT EXISTS user_vocabulary_progress (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    vocabulary_id INTEGER NOT NULL REFERENCES vocabulary(id) ON DELETE CASCADE,
    ease_factor DECIMAL(3,2) DEFAULT 2.5,  -- SM-2 algorithm ease factor
    interval INTEGER DEFAULT 1,             -- Days until next review
    repetitions INTEGER DEFAULT 0,          -- Number of successful reviews
    next_review_date TIMESTAMP WITH TIME ZONE NOT NULL,
    last_reviewed_at TIMESTAMP WITH TIME ZONE,
    total_reviews INTEGER DEFAULT 0,
    correct_reviews INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_user_vocabulary UNIQUE(user_id, vocabulary_id)
);

-- Create indexes for user vocabulary progress
CREATE INDEX idx_uvp_user_id ON user_vocabulary_progress(user_id);
CREATE INDEX idx_uvp_vocabulary_id ON user_vocabulary_progress(vocabulary_id);
CREATE INDEX idx_uvp_next_review ON user_vocabulary_progress(user_id, next_review_date);

-- Insert migration record
INSERT INTO schema_migrations (version) VALUES ('002_create_vocabulary_tables')
ON CONFLICT (version) DO NOTHING;
