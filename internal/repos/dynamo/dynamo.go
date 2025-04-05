package dynamo

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	db "github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func GetClient(cfg aws.Config) *db.Client {
	return db.NewFromConfig(cfg)
}
