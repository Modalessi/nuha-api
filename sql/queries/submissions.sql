-- name: CreateSubmission :one
INSERT INTO submissions (
    problem_id,
    user_id,
    language,
    source_code,
    status
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
) RETURNING *;


-- name: GetSubmissions :many
SELECT * FROM submissions ORDER BY created_at DESC OFFSET $1 LIMIT $2;


-- name: GetSubmissionByID :one
SELECT * FROM submissions WHERE id = $1;

-- name: GetSubmissionsByUserID :many
SELECT * FROM submissions WHERE user_id = $1 ORDER BY created_at DESC OFFSET $2 LIMIT $3;

-- name: GetSubmissionsByProblemID :many
SELECT * FROM submissions WHERE problem_id = $1 ORDER BY created_at DESC OFFSET $2 LIMIT $3;

-- name: GetUserSubmissionsForProblem :many
SELECT * FROM submissions 
WHERE user_id = $1 AND problem_id = $2 
ORDER BY created_at DESC
OFFSET $3 LIMIT $4;

-- name: UpdateSubmissionStatus :one
UPDATE submissions SET
    status = $2,
    updated_at = now()
WHERE id = $1 RETURNING *;

-- name: CreateSubmissionResult :one
INSERT INTO submission_results (
    id,
    submission_id,
    judge_token,
    stdin,
    stdout,
    expected_output,
    status_id,
    time_used,
    memory_used,
    judge_response
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10
) RETURNING *;


-- name: CreateSubmissionResults :many
INSERT INTO submission_results (
    id,
    submission_id,
    judge_token,
    stdin,
    stdout,
    expected_output,
    status_id,
    time_used,
    memory_used,
    judge_response
) 
VALUES (
    gen_random_uuid(),
    $1,
    unnest(@tokens::text[]),    
    unnest(@stdins::text[]),    
    unnest(@stdouts::text[]),   
    unnest(@expectedOutputs::text[]),
    unnest(@statuses::integer[]),
    unnest(@times::text[]),
    unnest(@memories::float8[]),
    unnest(@responses::bytea[])
)
RETURNING *;


-- name: GetSubmissionResultByID :one
SELECT * FROM submission_results WHERE id = $1;

-- name: GetSubmissionResultsBySubmissionID :many
SELECT * FROM submission_results WHERE submission_id = $1 OFFSET $2 LIMIT $3;


-- name: UpdateSubmissionResult :one
UPDATE submission_results SET
    stdin = $2,
    stdout = $3,
    expected_output = $4,
    status_id = $5,
    time_used = $6,
    memory_used = $7,
    judge_response = $8,
    updated_at = now()
WHERE id = $1 RETURNING *;

-- name: UpdateSubmissionResultStatus :one
UPDATE submission_results SET
    status_id = $2,
    updated_at = now()
WHERE id = $1 RETURNING *;