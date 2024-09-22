package wsconn

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/LaughG33k/chatWSService/iternal/client/redis"
	"github.com/LaughG33k/chatWSService/pkg"

	"github.com/gorilla/websocket"
)

type WsConnection struct {
	ctx         context.Context
	conn        *websocket.Conn
	redisClient *redis.RedisClient
	receiveChan chan []byte
	stopChan    chan struct{}
	connUuid    string
	connTimeout int64
	onceClose   sync.Once
	wp          *pkg.WorkerPool
	mu          sync.Mutex
}

func NewWsConn(ctx context.Context, conn *websocket.Conn, redisClient *redis.RedisClient, connUuid string, connTimeout int64) *WsConnection {

	return &WsConnection{

		ctx:         ctx,
		conn:        conn,
		redisClient: redisClient,
		receiveChan: make(chan []byte, 100),
		stopChan:    make(chan struct{}),
		connUuid:    connUuid,
		connTimeout: connTimeout,
		onceClose:   sync.Once{},
		wp:          pkg.NewWorkerPool(10),
	}

}

func (c *WsConnection) Start() {

	pkg.C.Add(func(ctx context.Context) error {
		return c.dropConn()
	})

	if err := c.redisClient.SubscribeOnGetMessage(c.connUuid, c.receiveChan); err != nil {
		fmt.Println(err)
	}

	go c.read()
	go func() {
		for {

			select {

			case _, ok := <-c.stopChan:
				if !ok {
					return
				}

			case data := <-c.receiveChan:

				if err := c.proccessIncomingMessage(data); err != nil {
					fmt.Println(err)
				}

			case <-c.ctx.Done():
				c.dropConn()

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
