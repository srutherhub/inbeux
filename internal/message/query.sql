-- name: GetMessageByExternalId :one
SELECT * FROM messages WHERE message_id = $1;

-- name: CreateMessage :exec
INSERT INTO messages (user_id, message_id, source) VALUES ($1, $2, $3);