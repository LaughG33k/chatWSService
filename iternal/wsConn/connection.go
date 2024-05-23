package wsconn

import (
	"context"
	"fmt"
	"sync"

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
	notAuth     bool
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
		notAuth:     false,
		wp:          pkg.NewWorkerPool(10),
	}

}

func (c *WsConnection) Start() {

	go c.read()
	go c.readFromRedis()

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

		default:
			if err := c.checkTimelifeConn(); err != nil {
				fmt.Println(err)
			}

		}

	}

}
