package chat

import (
	"context"
	"encoding/json"
	"time"

	"github.com/LaughG33k/chatWSService/iternal/model"
)

func (s *ChatService) Send(ctx context.Context, body model.Body101) error {

	bytes, err := json.Marshal(map[string]any{
		"type": 201,
		"body": map[string]any{
			"text":      body.Text,
			"messageId": body.MessageId,
			"sender":    body.Sender,
		},
	})

	if err != nil {
		return err
	}

	if err := s.publisher.Publish(ctx, body.Receiver, bytes); err != nil {
		return err
	}

	if err := s.messageRepo.SaveMessage(ctx, model.MessageForSave{
		SenderUuid:   body.Sender,
		ReceiverUuid: body.Receiver,
		MessageId:    body.MessageId,
		Text:         body.Text,
		Time:         time.Now().Unix(),
	}); err != nil {
		return err
	}

	return nil
}
