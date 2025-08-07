-- schema untuk table 'question'
CREATE TABLE questions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    content TEXT NOT NULL,
    hit INTEGER NOT NULL DEFAULT 0, -- untuk load balancing/randomisasi
    explanation_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- untuk performa query berdasarkan hit dan updated_at
CREATE INDEX idx_questions_hit_updated_at ON questions (hit, updated_at);

-- schema untuk table 'question_options'
CREATE TABLE question_options (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    is_correct BOOLEAN NOT NULL DEFAULT FALSE
);