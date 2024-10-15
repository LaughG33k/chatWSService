package chat

import (
	"context"

	"github.com/LaughG33k/chatWSService/iternal/model"
)

func (s *ChatService) GetHistory(ctx context.Context, body model.Body103) (model.MessageHistory, error) {

	return s.messageRepo.GetHistory(ctx, body.WithWhom)

}
