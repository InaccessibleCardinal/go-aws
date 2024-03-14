package users

import (
	"context"
	"go-aws/internal/conf"
	"go-aws/internal/repos"
	"go-aws/internal/repos/dynamo"
)

func InitUsersDB(ctx context.Context) UserDbIface {
	idGenerator := repos.CreateId
	cfg := conf.GetConfig(ctx)
	client := dynamo.GetClient(cfg)

	return New(ctx, client, idGenerator)
}
