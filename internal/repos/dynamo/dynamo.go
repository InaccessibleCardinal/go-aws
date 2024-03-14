package dynamo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	db "github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func GetClient(cfg aws.Config) *db.Client {
	return db.NewFromConfig(cfg)
}

type DynamoIface interface {
	DeleteItem(context.Context, *db.DeleteItemInput, ...func(*db.Options)) (*db.DeleteItemOutput, error)
	GetItem(context.Context, *db.GetItemInput, ...func(*db.Options)) (*db.GetItemOutput, error)
	PutItem(context.Context, *db.PutItemInput, ...func(*db.Options)) (*db.PutItemOutput, error)
}
