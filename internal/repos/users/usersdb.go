package users

import (
	"context"

	"go-aws/internal/types"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type GenId func(string) string
type DynamoIface interface {
	DeleteItem(
		context.Context,
		*dynamodb.DeleteItemInput,
		...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
	GetItem(
		context.Context,
		*dynamodb.GetItemInput,
		...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	PutItem(
		context.Context,
		*dynamodb.PutItemInput,
		...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

type UserRepo struct {
	client      DynamoIface
	idGenerator GenId
	tableName   string
}

func New(client DynamoIface, genId GenId) *UserRepo {
	return &UserRepo{client: client, idGenerator: genId, tableName: os.Getenv("AWS_DYNAMO_TABLE")}
}

func (d *UserRepo) PutUser(ctx context.Context, user types.User) error {
	user.UserID = d.idGenerator("USER#")
	marshaledUser, err := UserToAttributeValueMap(user)
	if err != nil {
		return err
	}
	_, err = d.client.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      marshaledUser,
		TableName: &d.tableName,
	})
	return err
}

func (d *UserRepo) GetUser(ctx context.Context, userID string) (*types.User, error) {
	out, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
		Key: IDToAttributeValue(userID), TableName: &d.tableName,
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, nil
	}
	user, err := UserFromAttributeValueMap(out.Item)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (d *UserRepo) RemoveUser(ctx context.Context, userId string) error {
	_, err := d.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key: IDToAttributeValue(userId), TableName: &d.tableName})
	if err != nil {
		return err
	}
	return nil
}
