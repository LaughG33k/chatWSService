package service

import (
	"context"

	"github.com/LaughG33k/chatWSService/iternal/model"
)

type ChatService interface {
	Send(ctx context.Context, body model.Body101) error
	GetHistory(ctx context.Context, body model.Body103) (model.MessageHistory, error)
	Edit(ctx context.Context, body model.MessageForEdit) error
	Delete(ctx context.Context, body model.Body104) error
}
