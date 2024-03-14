package cognito

import (
	"context"
	"go-aws/internal/conf"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

func InitCognito() *CognitoService {
	ctx := context.Background()
	cfg := conf.GetConfig(ctx)

	client := cognitoidentityprovider.NewFromConfig(cfg)
	getCsrp := GetCsrp
	return New(ctx, client, getCsrp)
}
