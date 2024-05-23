package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/LaughG33k/chatWSService/iternal/client/redis"
	wsconn "github.com/LaughG33k/chatWSService/iternal/wsConn"
	"github.com/LaughG33k/chatWSService/pkg"

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

	jwt := strings.Split(r.Header.Get("Jwt"), ".")

	if len(jwt) < 3 {
		http.Error(w, "bad jwt", http.StatusBadRequest)
		return
	}

	bytes, err := base64.RawStdEncoding.DecodeString(jwt[1])

	if err != nil {
		fmt.Println(err)
		http.Error(w, "bad jwt", http.StatusBadRequest)
		return
	}

	var playload map[string]any

	if err := json.Unmarshal(bytes, &playload); err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	if !pkg.CheckForAllKeys(playload, "uuid", "exp") {
		http.Error(w, "bad jwt. required fields are missing", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Panic(err)
		return
	}

	wsConn := wsconn.NewWsConn(h.ctx, conn, h.redisClient, playload["uuid"].(string), int64(playload["exp"].(float64)))

	if wsConn == nil {
		return
	}

	go wsConn.Start()

}
