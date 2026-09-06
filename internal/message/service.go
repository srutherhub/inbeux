package message

import "context"

type MessageService struct {
	query *MessageRepository
}

func NewService(repository *MessageRepository) *MessageService {
	return &MessageService{query: repository}
}

func (us *MessageService) ReceiveEmail(ctx context.Context, userId int64, messageId string) error {
	err := us.query.CreateMessage(ctx, userId, messageId, "email")
	return err
}
