package wsconn

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/LaughG33k/chatWSService/iternal/model"
	"github.com/LaughG33k/chatWSService/pkg"

	"github.com/goccy/go-json"
	"github.com/mitchellh/mapstructure"
)

func (c *WsConnection) proccessIncomingMessage(data []byte) error {

	var message model.WsMessage

	if err := json.Unmarshal(data, &message); err != nil {
		return err
	}

	switch message.Type {

	case 107:

		if err := c.validateNewJwt(message.Body); err != nil {

			bytes, err := json.Marshal(map[string]any{
				"type": 209,
				"body": map[string]any{
					"success": false,
					"error":   err.Error(),
				},
			})

			if err != nil {
				return err
			}

			c.send(bytes)

			break

		}

		bytes, err := json.Marshal(map[string]any{
			"type": 209,
			"body": map[string]any{
				"success": true,
				"error":   "",
			},
		})

		if err != nil {
			return err
		}

		c.send(bytes)

	case 101:

		if c.notAuth {
			return fmt.Errorf("access denied")
		}

		if err := c.sendMessage(message.Body); err != nil {
			return err
		}

	case 201:

		if c.notAuth {
			return fmt.Errorf("access denied")
		}
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

func (c *WsConnection) validateNewJwt(content map[string]any) error {

	if val, ok := content["jwt"].(string); !ok {
		return fmt.Errorf("Incorrect information, please fill in all fields")
	} else if val == "" {
		return fmt.Errorf("Incorrect information, please fill in all fields")
	}

	jwt := strings.Split((content["jwt"].(string)), ".")

	if len(jwt) < 3 {
		return fmt.Errorf("bad jwt")
	}

	bytes, err := base64.RawStdEncoding.DecodeString(jwt[1])

	if err != nil {
		return fmt.Errorf("bad jwt")
	}

	var playload map[string]any

	if err := json.Unmarshal(bytes, &playload); err != nil {
		return fmt.Errorf("can not unmarshal")
	}

	if !pkg.CheckForAllKeys(playload, "uuid", "exp") {
		return fmt.Errorf("bad jwt")
	}

	if playload["uuid"].(string) != c.connUuid {
		return fmt.Errorf("does not match the previous uuid")
	}

	c.connUuid = playload["uuid"].(string)

	if int64(playload["exp"].(float64)) <= time.Now().Unix() {
		return fmt.Errorf("token has expired")
	}

	c.connTimeout = int64(playload["exp"].(float64))

	c.notAuth = false

	return nil

}
