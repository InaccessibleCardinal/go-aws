package server

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	db "github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func GetClient(cfg aws.Config) *db.Client {
	return db.NewFromConfig(cfg)
}

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
