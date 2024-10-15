package chat

import (
	"context"

	"github.com/LaughG33k/chatWSService/iternal/model"
	"github.com/goccy/go-json"
)

func (s *ChatService) Delete(ctx context.Context, body model.Body104) error {

	if body.FlagDelForEvr {
		if err := s.messageRepo.DelMsgForEvryone(ctx, model.MessageForDelete{
			Sender:    body.Sender,
			Receiver:  body.WithWhom,
			MessageId: body.MessageId,
		}); err != nil {

			bytes, err := json.Marshal(map[string]any{
				"type": 202,
				"body": map[string]any{
					"withWhom":  body.Sender,
					"messageId": body.MessageId,
				},
			})

			if err != nil {
				return err
			}

			if err := s.publisher.Publish(ctx, body.WithWhom, bytes); err != nil {
				return err
			}

			return err
		}
	} else {
		if err := s.messageRepo.DeleteMessage(ctx, model.MessageForDelete{
			Sender:    body.Sender,
			Receiver:  body.WithWhom,
			MessageId: body.MessageId,
		}); err != nil {
			return err
		}
	}

	return nil
}
