package message

import (
	"context"
	"fmt"
	messagedb "inbeux/internal/message/db"

	"github.com/jackc/pgx/v5"
)

type MessageRepository struct {
	conn    *pgx.Conn
	queries *messagedb.Queries
}

func NewRepository(db *pgx.Conn) *MessageRepository {
	return &MessageRepository{conn: db, queries: messagedb.New(db)}
}

func (mr *MessageRepository) CreateMessage(ctx context.Context, userId int64, messageId string, source messagedb.MessagesSource) error {
	if userId == 0 {
		return fmt.Errorf("invalid argument: userId cannot be null")
	}

	if messageId == "" {
		return fmt.Errorf("invalid argument: messageId cannot be null")
	}

	if source == "" {
		return fmt.Errorf("invalid argument: source cannot be null")
	}

	input := messagedb.CreateMessageParams{UserID: userId, MessageID: messageId, Source: source}

	err := mr.queries.CreateMessage(ctx, input)

	return err
}

func (mr *MessageRepository) GetPendingMessages(ctx context.Context) ([]messagedb.Message, error) {

	var NUM_MESSAGES int32 = 50
	messages, err := mr.queries.GetPendingMessages(ctx, NUM_MESSAGES)

	if err != nil {
		return []messagedb.Message{}, nil
	}

	return messages, nil
}
