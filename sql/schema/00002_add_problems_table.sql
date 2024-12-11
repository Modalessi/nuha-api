-- +goose Up
-- +goose StatementBegin
CREATE TABLE problems (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    difficulty VARCHAR(255) NOT NULL,
    tags TEXT[] NOT NULL,
    time_limit FLOAT NOT NULL DEFAULT 1,
    memory_limit FLOAT NOT NULL DEFAULT 128000,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE problems_descriptions (
    problem_id UUID PRIMARY KEY REFERENCES problems(id),
    description TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE test_cases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    problem_id UUID NOT NULL REFERENCES problems(id),
    number INTEGER NOT NULL,
    stdin TEXT NOT NULL,
    expected_output TEXT NOT NULL,
    
    UNIQUE(problem_id, number)
);
CREATE INDEX idx_test_cases_problem_id ON test_cases(problem_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE problems_descriptions;
DROP TABLE test_cases;
DROP TABLE problems;
-- +goose StatementEnd
