-- name: GetActionsForMessage :many
SELECT * FROM actions WHERE message_id = $1;