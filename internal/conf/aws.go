package conf

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func GetConfig(ctx context.Context) aws.Config {
	conf, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("failed to create config %s\n", err)
	}
	return conf
}
