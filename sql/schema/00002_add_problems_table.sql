-- +goose Up
-- +goose StatementBegin
CREATE TABLE problems (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    difficulty VARCHAR(255) NOT NULL,
    description_path VARCHAR(255) NOT NULL,
    testcases_path VARCHAR(255) NOT NULL,
    tags TEXT[] NOT NULL,
    time_limit FLOAT NOT NULL DEFAULT 1,
    memory_limit FLOAT NOT NULL DEFAULT 128000,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE problems;
-- +goose StatementEnd
