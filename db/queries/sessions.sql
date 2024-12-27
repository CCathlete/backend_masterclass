-- name: CreateSession :one
INSERT INTO sessions (
  id,
  username,
  refresh_token,
  user_agent,
  client_ip,
  is_blocked,
  expires_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions
WHERE id = $1 LIMIT 1;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = $1;

-- name: UpdateSession :one
UPDATE sessions
SET
  refresh_token = $2,
  user_agent = $3,
  client_ip = $4,
  is_blocked = $5,
  expires_at = $6
WHERE
  id = $1
RETURNING *;

-- name: ListSessions :many
SELECT * FROM sessions
WHERE username = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: CountSessions :one
SELECT COUNT(*) FROM sessions
WHERE username = $1;

-- name: ListBlockedSessions :many
SELECT * FROM sessions
WHERE is_blocked = true
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: CountBlockedSessions :one
SELECT COUNT(*) FROM sessions
WHERE is_blocked = true;

-- name: BlockSession :exec
UPDATE sessions
SET is_blocked = true
WHERE id = $1;

-- name: UnblockSession :exec
UPDATE sessions
SET is_blocked = false
WHERE id = $1;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions
WHERE expires_at < NOW();

-- name: DeleteBlockedSessions :exec
DELETE FROM sessions
WHERE is_blocked = true;

