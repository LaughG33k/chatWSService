package wsconn

import (
	"context"
	"sync"
	"time"

	"github.com/LaughG33k/chatWSService/iternal/client/redis"
	"github.com/LaughG33k/chatWSService/iternal/repository"
	"github.com/LaughG33k/chatWSService/pkg"

	"github.com/gorilla/websocket"
)

type WsConnection struct {
	conn        *websocket.Conn
	redisClient *redis.RedisClient
	msgRepo     repository.Messages
	receiveChan chan []byte
	stopChan    chan struct{}
	connUuid    string
	onceClose   sync.Once
	wp          *pkg.WorkerPool
	mu          sync.Mutex
}

func NewWsConn(conn *websocket.Conn, redisClient *redis.RedisClient, msgRepo repository.Messages, connUuid string) *WsConnection {

	return &WsConnection{
		conn:        conn,
		redisClient: redisClient,
		msgRepo:     msgRepo,
		receiveChan: make(chan []byte, 100),
		stopChan:    make(chan struct{}),
		connUuid:    connUuid,
		onceClose:   sync.Once{},
		wp:          pkg.NewWorkerPool(10),
	}

}

func (c *WsConnection) Start() {

	if err := c.redisClient.SubscribeOnGetMessage(c.connUuid, c.receiveChan); err != nil {
		pkg.Log.Infof("error with redis %s: %s", c.connUuid, err)
	}

	go c.read()
	go func() {
		defer func() {
			pkg.Log.Infof("exit from process gorutine: %s", c.connUuid)
		}()
		for {

			select {

			case _, ok := <-c.stopChan:
				if !ok {
					return
				}

			case data := <-c.receiveChan:

				c.wp.AddWorker(func() {
					tm, canc := context.WithTimeout(context.TODO(), 45*time.Second)
					defer canc()
					pkg.Log.Infof("start worker for %s", c.connUuid)
					if err := c.proccessIncomingMessage(tm, data); err != nil {
						pkg.Log.Infof("error for conn %s: %s", c.connUuid, err)
					}
					pkg.Log.Infof("finished worker for %s", c.connUuid)
				})

			}

		}

	}()

}

func (c *WsConnection) Close() error {
	return c.dropConn()
}

func (c *WsConnection) dropConn() (err error) {

	c.onceClose.Do(func() {
		close(c.stopChan)
		time.Sleep(1 * time.Minute)
		err = c.conn.Close()
	})

	return err

}
