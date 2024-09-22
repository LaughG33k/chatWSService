package repository

import (
	"context"

	"github.com/LaughG33k/chatWSService/iternal/model"
)

type Messages interface {
	SaveMessage(context.Context, model.MessageForSave) error
	DeleteMessage(context.Context, model.MessageForDelete) error
	DelMsgForEvryone(context.Context, model.MessageForDelete) error
	EditMessage(context.Context, model.MessageForEdit) error
	GetHistory(context.Context, string) (model.MessageHistory, error)
}
