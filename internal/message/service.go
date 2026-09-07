package message

import (
	"context"
	"fmt"
	messagedb "inbeux/internal/message/db"
	"inbeux/internal/utils"
	"time"
)

type MessageService struct {
	query *MessageRepository
}

func NewService(repository *MessageRepository) *MessageService {
	return &MessageService{query: repository}
}

func (ms *MessageService) ReceiveEmail(ctx context.Context, userId int64, messageId string) error {
	err := ms.query.CreateMessage(ctx, userId, messageId, "email")
	return err
}

func (ms *MessageService) StartPendingMessagesBatch(ctx context.Context) {
	workerDefinition := &pendingMessagesWorker{messageService: ms}
	worker := utils.NewWorker(workerDefinition, 25)
	worker.Start(ctx, 5*time.Second)
}

type pendingMessagesWorker struct {
	messageService *MessageService
}

func (esw *pendingMessagesWorker) GetTasks(ctx context.Context) ([]messagedb.Message, error) {
	messages, err := esw.messageService.query.GetPendingMessages(ctx)

	if err != nil {
		return []messagedb.Message{}, fmt.Errorf("failed to retrieve pending message: %w", err)
	}

	return messages, nil
}

func (esw *pendingMessagesWorker) Process(ctx context.Context, message messagedb.Message) error {
	fmt.Println(message)
	return nil
}
