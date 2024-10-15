package chat

import (
	"github.com/LaughG33k/chatWSService/iternal/queue"
	"github.com/LaughG33k/chatWSService/iternal/repository"
	"github.com/LaughG33k/chatWSService/iternal/service"
)

type ChatService struct {
	messageRepo repository.Messages
	publisher   queue.Publisher
}

func NewChatService(messageRepo repository.Messages, publisher queue.Publisher) service.ChatService {
	return &ChatService{
		messageRepo: messageRepo,
		publisher:   publisher,
	}
}
