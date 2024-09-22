package handler

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/LaughG33k/chatWSService/iternal/client/redis"
	wsconn "github.com/LaughG33k/chatWSService/iternal/wsConn"

	"github.com/gorilla/websocket"
)

var upgrader *websocket.Upgrader = &websocket.Upgrader{
	ReadBufferSize:   0,
	WriteBufferSize:  0,
	WriteBufferPool:  &sync.Pool{},
	HandshakeTimeout: 5 * time.Minute,

	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WsChatHandler struct {
	ctx         context.Context
	redisClient *redis.RedisClient
}

func NewWsChatHandler(ctx context.Context, redisClient *redis.RedisClient) *WsChatHandler {
	return &WsChatHandler{
		ctx:         ctx,
		redisClient: redisClient,
	}
}

func (h *WsChatHandler) StartHandler() {

	http.HandleFunc("/ws/chat", h.startWsConn)

}

func (h *WsChatHandler) startWsConn(w http.ResponseWriter, r *http.Request) {

	uuid := r.Header.Get("User-Uuid")

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Panic(err)
		return
	}

	wsConn := wsconn.NewWsConn(h.ctx, conn, h.redisClient, uuid, 0)

	if wsConn == nil {
		return
	}

	wsConn.Start()

}
