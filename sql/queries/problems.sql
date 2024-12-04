-- name: GetProblems :many
SELECT * FROM problems OFFSET $1 LIMIT $2;


-- name: GetProblemByID :one
SELECT * FROM problems WHERE id = $1;


-- name: CreateProblem :one
INSERT INTO problems (
    id,
    title,
    difficulty,
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
    $7,
    $8
) RETURNING *;
