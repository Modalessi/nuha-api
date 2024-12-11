-- name: GetProblems :many
SELECT * FROM problems OFFSET $1 LIMIT $2;


-- name: GetProblemByID :one
SELECT * FROM problems WHERE id = $1;


-- name: CreateProblem :one
INSERT INTO problems (
    title,
    difficulty,
    tags,
    time_limit,
    memory_limit
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
) RETURNING *;


-- name: DeleteProblem :one
DELETE FROM problems WHERE id = $1 RETURNING *;


-- name: UpdateProblem :one
UPDATE problems SET
    title = $2,
    difficulty = $3,
    tags = $4,
    time_limit = $5,
    memory_limit = $6,
    updated_at = now()
WHERE id = $1 RETURNING *;


-- name: AddProblemDescription :one
INSERT INTO problems_descriptions (
    problem_id,
    description
) VALUES (
    $1,
    $2
) RETURNING *;


-- name: GetProblemDescription :one
SELECT * FROM problems_descriptions WHERE problem_id = $1;


-- name: DeleteProblemDescription :one
DELETE FROM problems_descriptions WHERE problem_id = $1 RETURNING *;

-- name: UpdateProblemDescription :one
UPDATE problems_descriptions SET
    description = $2,
    updated_at = now()
WHERE problem_id = $1 RETURNING *;


-- name: CreateTestCases :many
WITH numbered_arrays AS (
    SELECT 
        generate_series(
            COALESCE((SELECT MAX(number) FROM test_cases WHERE problem_id = $1), 0) + 1,
            COALESCE((SELECT MAX(number) FROM test_cases WHERE problem_id = $1), 0) + array_length(@stdins::TEXT[], 1)
        ) as num,
        unnest(@stdins::TEXT[]) as in_data,
        unnest(@expected_outputs::TEXT[]) as out_data
)
INSERT INTO test_cases (
    problem_id,
    number,
    stdin,
    expected_output
) 
SELECT 
    $1,
    num,
    in_data,
    out_data
FROM numbered_arrays
RETURNING *;


-- name: GetTestCases :many
SELECT * FROM test_cases WHERE problem_id = $1 ORDER BY number;

-- name: DeleteTestCases :many
DELETE FROM test_cases WHERE problem_id = $1 RETURNING *;