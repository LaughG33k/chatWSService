package wsconn

import (
	"context"

	"github.com/LaughG33k/chatWSService/iternal/model"
	"github.com/goccy/go-json"
	"github.com/mitchellh/mapstructure"
)

func (c *WsConnection) handle(ctx context.Context, data []byte) error {

	message := &model.WsMessage{}

	if err := json.Unmarshal(data, message); err != nil {
		return err
	}

	switch message.Type {

	case 101:
		if err := c.messageForSend(ctx, message.Body); err != nil {
			return err
		}
	}

}

func (c *WsConnection) messageForSend(ctx context.Context, data map[string]any) error {

	body := model.Body101{}

	if err := mapstructure.Decode(data, &body); err != nil {
		return err
	}

	if err := c.chatService.Send(ctx, body); err != nil {
		return err
	}

	return nil
}

func (c *WsConnection) messageForDelete(ctx context.Context, data map[string]any) error {

	return nil
}
