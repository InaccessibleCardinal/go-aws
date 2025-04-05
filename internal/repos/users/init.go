package users

import (
	"context"
	"go-aws/internal/conf"
	"go-aws/internal/repos/dynamo"
	"go-aws/internal/repos/ids"
)

func InitUsersDB(ctx context.Context) *UserRepo {
	idGenerator := ids.CreateId
	cfg := conf.GetConfig(ctx)
	client := dynamo.GetClient(cfg)

	return New(client, idGenerator)
}
