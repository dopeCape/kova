-- name: CreateAccount :one
INSERT INTO accounts (user_id, github_username, github_id, avatar_url, access_token)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, github_username, github_id, avatar_url, created_at, updated_at;

-- name: GetAccountByID :one
SELECT id, user_id, github_username, github_id, avatar_url, access_token, created_at, updated_at
FROM accounts
WHERE id = $1;

-- name: GetAccountByGithubID :one
SELECT id, user_id, github_username, github_id, avatar_url, access_token, created_at, updated_at
FROM accounts
WHERE github_id = $1;

-- name: GetAccountByGithubUsername :one
SELECT id, user_id, github_username, github_id, avatar_url, access_token, created_at, updated_at
FROM accounts
WHERE github_username = $1;

-- name: GetAccountsByUserID :many
SELECT id, user_id, github_username, github_id, avatar_url, created_at, updated_at
FROM accounts
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetAccountsByUserIDWithTokens :many
SELECT id, user_id, github_username, github_id, avatar_url, access_token, created_at, updated_at
FROM accounts
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateAccount :one
UPDATE accounts
SET github_username = $2, avatar_url = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, user_id, github_username, github_id, avatar_url, created_at, updated_at;

-- name: UpdateAccountToken :one
UPDATE accounts
SET access_token = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, user_id, github_username, github_id, avatar_url, created_at, updated_at;

-- name: UpdateAccountByGithubID :one
UPDATE accounts
SET github_username = $2, avatar_url = $3, access_token = $4, updated_at = CURRENT_TIMESTAMP
WHERE github_id = $1
RETURNING id, user_id, github_username, github_id, avatar_url, created_at, updated_at;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1;

-- name: DeleteAccountByGithubID :exec
DELETE FROM accounts
WHERE github_id = $1;

-- name: DeleteAccountsByUserID :exec
DELETE FROM accounts
WHERE user_id = $1;

-- name: ListAccounts :many
SELECT id, user_id, github_username, github_id, avatar_url, created_at, updated_at
FROM accounts
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountAccounts :one
SELECT COUNT(*) FROM accounts;

-- name: CountAccountsByUserID :one
SELECT COUNT(*) FROM accounts WHERE user_id = $1;

-- name: AccountExistsByGithubID :one
SELECT EXISTS(SELECT 1 FROM accounts WHERE github_id = $1);

-- name: AccountExistsByGithubUsername :one
SELECT EXISTS(SELECT 1 FROM accounts WHERE github_username = $1);

-- name: AccountExistsByUserIDAndGithubID :one
SELECT EXISTS(SELECT 1 FROM accounts WHERE user_id = $1 AND github_id = $2);

-- name: AccountExistsForUser :one
SELECT EXISTS(SELECT 1 FROM accounts WHERE user_id = $1 AND github_id = $2);

-- name: SearchAccounts :many
SELECT id, user_id, github_username, github_id, avatar_url, created_at, updated_at
FROM accounts
WHERE github_username ILIKE '%' || $1 || '%'
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: SearchAccountsByUserID :many
SELECT id, user_id, github_username, github_id, avatar_url, created_at, updated_at
FROM accounts
WHERE user_id = $1 AND github_username ILIKE '%' || $2 || '%'
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: GetAccountWithUser :one
SELECT 
    a.id,
    a.user_id,
    a.github_username,
    a.github_id,
    a.avatar_url,
    a.access_token,
    a.created_at,
    a.updated_at,
    u.username,
    u.email
FROM accounts a
JOIN users u ON a.user_id = u.id
WHERE a.id = $1;


