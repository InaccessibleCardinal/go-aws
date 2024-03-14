package users

import (
	"context"
	"encoding/json"
	"go-aws/internal/repos"
	"go-aws/internal/repos/dynamo"
	"go-aws/internal/types"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type UserDbIface interface {
	GetUser(userId string) (*types.UserRecord, error)
	PutUser(rawUser types.UserRecord) error
	RemoveUser(userId string) error
}

type UserDb struct {
	client      dynamo.DynamoIface
	ctx         context.Context
	idGenerator repos.GenId
	tableName   string
}

func (d *UserDb) PutUser(user types.UserRecord) error {
	user.UserId = d.idGenerator("USER#")
	marshaledUser, err := userToAttributeValueMap(user)
	if err != nil {
		return err
	}
	_, err = d.client.PutItem(d.ctx, &dynamodb.PutItemInput{
		Item:      marshaledUser,
		TableName: &d.tableName,
	})
	return err
}

func (d *UserDb) GetUser(userId string) (*types.UserRecord, error) {
	out, err := d.client.GetItem(d.ctx, &dynamodb.GetItemInput{Key: idToAttributeValue(userId), TableName: &d.tableName})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, nil
	}
	user, err := userFromAttributeValueMap(out.Item)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (d *UserDb) RemoveUser(userId string) error {
	_, err := d.client.DeleteItem(d.ctx, &dynamodb.DeleteItemInput{
		Key: idToAttributeValue(userId), TableName: &d.tableName})
	if err != nil {
		return err
	}
	return nil
}

func New(ctx context.Context, client dynamo.DynamoIface, genId repos.GenId) UserDbIface {
	return &UserDb{ctx: ctx, client: client, idGenerator: genId, tableName: os.Getenv("AWS_DYNAMO_TABLE")}
}

func SaveJson(data any) {
	bts, _ := json.Marshal(data)
	os.WriteFile("response.json", bts, 0777)
}
