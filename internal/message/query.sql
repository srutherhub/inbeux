-- name: GetMessageByExternalId :one
SELECT * FROM messages WHERE message_id = $1;

-- name: CreateMessage :exec
INSERT INTO messages (user_id, message_id, source) VALUES ($1, $2, $3) ON CONFLICT (message_id) DO NOTHING;

-- name: GetPendingMessages :many
SELECT * FROM messages WHERE status = 'pending' ORDER BY created_at ASC LIMIT $1 FOR UPDATE SKIP LOCKED;

-- name: UpdateMessageToClassified :exec
UPDATE messages SET status = 'pending' WHERE message_id = $1;