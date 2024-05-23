package wsconn

import (
	"time"

	"github.com/goccy/go-json"
)

func (c *WsConnection) dropConn() {

	select {

	case _, ok := <-c.stopChan:
		if !ok {
			return
		}

	default:
		close(c.stopChan)
		c.conn.Close()
	}

}

func (c *WsConnection) checkTimelifeConn() error {

	if time.Now().Unix() > c.connTimeout {

		if !c.notAuth {

			c.notAuth = true

			go func() {
				time.Sleep(1 * time.Minute)
				c.dropConn()
			}()

			bytes, err := json.Marshal(map[string]any{
				"type": 210,
			})

			if err != nil {
				return err
			}

			c.send(bytes)

		}

	}

	return nil

}
