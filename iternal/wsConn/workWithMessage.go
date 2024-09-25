package wsconn

import (
	"context"
	"fmt"
	"time"

	"github.com/LaughG33k/chatWSService/iternal/model"

	"github.com/goccy/go-json"
	"github.com/mitchellh/mapstructure"
)

func (c *WsConnection) proccessIncomingMessage(ctx context.Context, data []byte) error {

	var message model.WsMessage

	if err := json.Unmarshal(data, &message); err != nil {
		return err
	}

	switch message.Type {

	case 101:

		if err := c.sendMessage(ctx, message.Body); err != nil {
			return err
		}

	case 201:

		c.send(data)

	}

	return nil
}

func (c *WsConnection) sendMessage(ctx context.Context, content map[string]any) error {

	body := &model.Body101{}

	if err := mapstructure.Decode(content, body); err != nil {
		return err
	}

	if body.Text == "" || body.Receiver == "" || body.MessageId == "" {
		return fmt.Errorf("uncorrect input data")
	}

	bytes, err := json.Marshal(map[string]any{
		"type": 201,
		"body": map[string]any{
			"text":      body.Text,
			"messageId": body.MessageId,
			"sender":    c.connUuid,
		},
	})

	if err != nil {
		return err
	}

	err = c.redisClient.PublishMessageToSend(body.Receiver, bytes)

	if err != nil {
		return err
	}

	if err := c.msgRepo.SaveMessage(ctx, model.MessageForSave{
		SenderUuid:   c.connUuid,
		ReceiverUuid: body.Receiver,
		Text:         body.Text,
		MessageId:    body.MessageId,
		Time:         time.Now().Unix(),
	}); err != nil {
		return err
	}

	bytes, err = json.Marshal(map[string]any{
		"messageId": body.MessageId,
		"success":   true,
	})

	if err != nil {
		return err
	}

	c.send(bytes)

	return nil

}

func (c *WsConnection) getHistory(ctx context.Context, content map[string]any) error {

	body := model.Body103{}

	if err := mapstructure.Decode(content, &body); err != nil {
		return err
	}

	history, err := c.msgRepo.GetHistory(ctx, body.WithWhom)

	if err != nil {
		return err
	}

	bytes, err := json.Marshal(map[string]any{
		"type": 206,
		"body": history,
	})

	if err != nil {
		return err
	}

	if err := c.send(bytes); err != nil {
		return err
	}

	return nil
}

func (c *WsConnection) editMsg(ctx context.Context, content map[string]any) error {

	body := model.Body105{}

	if err := mapstructure.Decode(content, &body); err != nil {
		return err
	}

	if err := c.msgRepo.EditMessage(ctx, model.MessageForEdit{
		Sender:    c.connUuid,
		Recipient: body.WithWhom,
		MessageId: body.MessageId,
		NewText:   body.UpdatedText,
	}); err != nil {
		return err
	}

	bytes, err := json.Marshal(map[string]any{
		"type": 208,
		"body": map[string]any{
			"messageId": body.MessageId,
			"success":   true,
		},
	})

	if err != nil {
		return err
	}

	if err := c.send(bytes); err != nil {
		return err
	}

	return nil
}

func (c *WsConnection) delMsg(ctx context.Context, content map[string]any) error {

	body := model.Body104{}

	if err := mapstructure.Decode(content, &body); err != nil {
		return err
	}

	if body.FlagDelForEvr {
		if err := c.msgRepo.DelMsgForEvryone(ctx, model.MessageForDelete{
			Sender:    c.connUuid,
			Receiver:  body.WithWhom,
			MessageId: body.MessageId,
		}); err != nil {
			return err
		}
	} else {
		if err := c.msgRepo.DeleteMessage(ctx, model.MessageForDelete{
			Sender:    c.connUuid,
			Receiver:  body.WithWhom,
			MessageId: body.MessageId,
		}); err != nil {
			return err
		}
	}

	bytes, err := json.Marshal(map[string]any{
		"type": 207,
		"body": map[string]any{
			"messageId": body.MessageId,
			"success":   true,
		},
	})

	if err != nil {
		return err
	}

	if err := c.send(bytes); err != nil {
		return err
	}

	return nil
}
