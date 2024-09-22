package messagesrepository

import (
	"context"
	"errors"
	"fmt"

	mcl "github.com/LaughG33k/chatWSService/iternal/client/mongo"
	"github.com/LaughG33k/chatWSService/iternal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatRepository struct {
	mongoClient *mcl.MongoClient
	collection  string
}

func NewRepository(mongoClient *mcl.MongoClient, collection string) *ChatRepository {

	return &ChatRepository{

		mongoClient: mongoClient,
		collection:  collection,
	}

}

func (c *ChatRepository) SaveMessage(ctx context.Context, data model.MessageForSave) error {

	updateAtRecipient := mongo.NewUpdateOneModel()
	updateAtRecipient.
		SetFilter(bson.M{"uuid": data.ReceiverUuid}).
		SetUpdate(bson.M{"$set": bson.M{fmt.Sprintf("received.%s.%s", data.SenderUuid, data.MessageId): bson.M{"text": data.Text, "time": data.Time, "editFlag": false}}})

	updateAtSender := mongo.NewUpdateOneModel()

	updateAtSender.
		SetFilter(bson.M{"uuid": data.SenderUuid}).
		SetUpdate(bson.M{"$set": bson.M{fmt.Sprintf("sent.%s.%s", data.ReceiverUuid, data.MessageId): bson.M{"text": data.Text, "time": data.Time, "editFlag": false}}})

	if _, err := c.mongoClient.Collection(c.collection).BulkWrite(ctx, []mongo.WriteModel{updateAtRecipient, updateAtSender}); err != nil {
		return err
	}

	return nil

}

func (c *ChatRepository) DeleteMessage(ctx context.Context, msg model.MessageForDelete) error {

	if _, err := c.mongoClient.Collection(c.collection).UpdateOne(ctx, bson.M{"uuid": msg.Sender}, bson.M{"$unset": bson.M{fmt.Sprintf("sent.%s.%s", msg.Receiver, msg.MessageId): 1}}); err != nil {
		return err
	}

	return nil

}

func (c *ChatRepository) DelMsgForEvryone(ctx context.Context, msg model.MessageForDelete) error {

	deleteAtReceiver := mongo.NewUpdateOneModel()

	deleteAtReceiver.SetFilter(bson.M{"uuid": msg.Receiver})
	deleteAtReceiver.SetUpdate(bson.M{"$unset": bson.M{fmt.Sprintf("received.%s.%s", msg.Sender, msg.MessageId): 1}})

	deleteAtSender := mongo.NewUpdateOneModel()

	deleteAtSender.SetFilter(bson.M{"uuid": msg.Sender})
	deleteAtSender.SetUpdate(bson.M{"$unset": bson.M{fmt.Sprintf("sent.%s.%s", msg.Receiver, msg.MessageId): 1}})

	if _, err := c.mongoClient.Collection(c.collection).BulkWrite(ctx, []mongo.WriteModel{deleteAtReceiver, deleteAtSender}); err != nil {
		return err
	}

	return nil
}

func (c *ChatRepository) EditMessage(ctx context.Context, msg model.MessageForEdit) error {

	updateAtRecipient := mongo.NewUpdateOneModel()
	updateAtRecipient.
		SetFilter(bson.M{"uuid": msg.Recipient}).
		SetUpdate(bson.M{"$set": bson.M{fmt.Sprintf("received.%s.%s.text", msg.Sender, msg.MessageId): msg.NewText, fmt.Sprintf("received.%s.%s.editFlag", msg.Sender, msg.MessageId): true}})
	updateAtSender := mongo.NewUpdateOneModel()

	updateAtSender.
		SetFilter(bson.M{"uuid": msg.Sender}).
		SetUpdate(bson.M{"$set": bson.M{fmt.Sprintf("sent.%s.%s.text", msg.Recipient, msg.MessageId): msg.NewText, fmt.Sprintf("sent.%s.%s.editFlag", msg.Recipient, msg.MessageId): true}})

	if _, err := c.mongoClient.Collection(c.collection).BulkWrite(ctx, []mongo.WriteModel{updateAtRecipient, updateAtSender}); err != nil {
		return err
	}

	return nil

}

func (c *ChatRepository) GetHistory(ctx context.Context, who string) (model.MessageHistory, error) {

	var history model.MessageHistory

	res := c.mongoClient.Collection(c.collection).FindOne(ctx, bson.M{"uuid": who})

	if res.Err() != nil {
		return history, res.Err()
	}

	if err := res.Decode(&history); err != nil {
		return history, errors.New("failed to decode history")
	}

	return history, nil

}
