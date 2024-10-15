package wsconn

import (
	"context"
	"sync"
	"time"

	"github.com/LaughG33k/chatWSService/iternal/queue"
	"github.com/LaughG33k/chatWSService/iternal/service"
	"github.com/LaughG33k/chatWSService/pkg"

	"github.com/gorilla/websocket"
)

type ReadWriter interface {
	Read() []byte
	Write() error
}

type WsConnection struct {
	conn        *websocket.Conn
	queue       queue.Queue
	receiveChan chan any
	chatService service.ChatService
	stopChan    chan struct{}
	connUuid    string
	onceClose   sync.Once
	wp          *pkg.WorkerPool
	mu          sync.Mutex
}

func NewWsConn(conn *websocket.Conn, chatService service.ChatService, queue queue.Queue, connUuid string) *WsConnection {

	return &WsConnection{
		conn:        conn,
		receiveChan: make(chan any, 100),
		stopChan:    make(chan struct{}),
		chatService: chatService,
		queue:       queue,
		connUuid:    connUuid,
		onceClose:   sync.Once{},
		wp:          pkg.NewWorkerPool(10),
	}

}

func (c *WsConnection) Start(ctx context.Context) {

	if err := c.queue.Subscribe(ctx, c.connUuid, c.receiveChan); err != nil {
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
