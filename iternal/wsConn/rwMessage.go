package wsconn

import (
	"fmt"

	"github.com/LaughG33k/chatWSService/pkg"
	"github.com/gorilla/websocket"
)

func (c *WsConnection) read() {

	defer func() {
		pkg.Log.Infof("exit from read gorutine: %s", c.connUuid)
	}()

	for {

		select {

		case _, ok := <-c.stopChan:
			if !ok {
				return
			}

		default:

			mt, data, err := c.conn.ReadMessage()

			if err != nil || mt == websocket.CloseMessage || mt == websocket.CloseAbnormalClosure {
				c.dropConn()
				fmt.Println(err)
				return
			}

			if len(data) == 0 {
				continue
			}

			c.receiveChan <- data

		}

	}

}

func (c *WsConnection) send(data []byte) error {

	select {

	case _, ok := <-c.stopChan:
		if !ok {
			return nil
		}

	default:
		c.mu.Lock()
		defer c.mu.Unlock()

		if err := c.conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
			return err
		}
	}

	return nil

}
