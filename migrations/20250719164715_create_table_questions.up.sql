-- schema untuk table 'question'
CREATE TABLE questions (
                           id CHAR(26) PRIMARY KEY,
                           content TEXT NOT NULL,
                           options JSONB DEFAULT '[]', -- pindahan dari question_options
                           hit INTEGER NOT NULL DEFAULT 0,      -- untuk load balancing/randomisasi
                           explanation_url TEXT,
                           created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                           updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- untuk performa query berdasarkan hit dan updated_at
CREATE INDEX idx_questions_hit_updated_at ON questions (hit, updated_at);