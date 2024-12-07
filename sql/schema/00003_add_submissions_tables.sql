-- +goose Up
-- +goose StatementBegin
CREATE TABLE submissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    problem_id UUID NOT NULL REFERENCES problems(id),
    user_id UUID NOT NULL REFERENCES users(id),
    language INTEGER NOT NULL,
    source_code TEXT NOT NULL,
    status VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_submissions_user_id ON submissions(user_id);
CREATE INDEX idx_submissions_problem_id ON submissions(problem_id);

CREATE TABLE submission_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    submission_id UUID NOT NULL REFERENCES submissions(id),
    judge_token VARCHAR(255) NOT NULL,
    stdin TEXT NOT NULL,
    stdout TEXT NOT NULL,
    expected_output TEXT NOT NULL,
    status_id INTEGER NOT NULL,
    time_used TEXT NOT NULL,
    memory_used FLOAT NOT NULL,
    judge_response bytea NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_submission_results_submission_id ON submission_results(submission_id);
CREATE INDEX idx_submission_results_judge_token ON submission_results(judge_token);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE submission_results;
DROP TABLE submissions;
-- +goose StatementEnd
