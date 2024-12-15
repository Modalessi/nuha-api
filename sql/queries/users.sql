-- name: CreateUser :one
INSERT INTO users (
    email,
    password
) VALUES (
    $1,
    $2
) RETURNING *;


-- name: SetUserVerified :one
UPDATE users SET verified = TRUE WHERE email = $1 RETURNING *;


-- name: CreateUserSession :one
INSERT INTO sessions (
    user_id,
    token,
    expires_at
) VALUES (
    $1,
    $2,
    $3
) RETURNING *;

-- name: CreateVerificationRequest :one
INSERT INTO verification_tokens (
    user_id,
    token
) VALUES (
    $1,
    $2
) RETURNING *;

-- name: CreatePasswordResetRequest :one
INSERT INTO password_reset_tokens (
    user_id,
    token
) VALUES (
    $1,
    $2
) RETURNING *;

-- name: SetSessionRevoked :one
UPDATE sessions SET revoked = TRUE WHERE token = $1 RETURNING *;;


-- name: GetVerficationToken :one
SELECT * FROM verification_tokens WHERE token = $1;


-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;


-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;