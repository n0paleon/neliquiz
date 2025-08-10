CREATE TABLE categories (
    id CHAR(26) PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL
);

-- m2m relation table
CREATE TABLE question_categories (
    question_id CHAR(26) REFERENCES questions(id) ON DELETE CASCADE,
    category_id CHAR(26) REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (question_id, category_id)
);