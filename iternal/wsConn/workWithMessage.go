package wsconn

import (
	"fmt"

	"github.com/LaughG33k/chatWSService/iternal/model"

	"github.com/goccy/go-json"
	"github.com/mitchellh/mapstructure"
)

func (c *WsConnection) proccessIncomingMessage(data []byte) error {

	var message model.WsMessage

	if err := json.Unmarshal(data, &message); err != nil {
		return err
	}

	switch message.Type {

	case 101:

		if err := c.sendMessage(message.Body); err != nil {
			return err
		}

	case 201:

		c.send(data)

	}

	return nil
}

func (c *WsConnection) sendMessage(content map[string]any) error {

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

	return nil

}
