-- name: GetProblems :many
SELECT (
    id,
    title,
    description_path,
    testcases_path,
    tags,
    time_limit,
    memory_limit,
    created_at,
    updated_at
) FROM problems;


-- name: GetProblemByID :one
SELECT (
    id,
    title,
    description_path,
    testcases_path,
    tags,
    time_limit,
    memory_limit,
    created_at,
    updated_at
) FROM problems WHERE id = $1;


-- name: CreateProblem :one
INSERT INTO problems (
    id,
    title,
    description_path,
    testcases_path,
    tags,
    time_limit,
    memory_limit
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
) RETURNING *;
