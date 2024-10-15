package chat

import (
	"context"

	"github.com/LaughG33k/chatWSService/iternal/model"
)

func (s *ChatService) Edit(ctx context.Context, body model.MessageForEdit) error {
	return s.messageRepo.EditMessage(ctx, body)
}
