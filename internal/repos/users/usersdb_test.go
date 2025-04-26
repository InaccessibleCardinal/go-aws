package users

import (
	"context"
	"errors"
	"go-aws/internal/types"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/assert"
)

type testDynamoDB struct {
	getOutput    *dynamodb.GetItemOutput
	errGet       error
	putOutput    *dynamodb.PutItemOutput
	errPut       error
	deleteOutput *dynamodb.DeleteItemOutput
	errDelete    error
}

func (td *testDynamoDB) DeleteItem(
	context.Context,
	*dynamodb.DeleteItemInput,
	...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	return td.deleteOutput, td.errDelete
}

func (td *testDynamoDB) GetItem(
	context.Context,
	*dynamodb.GetItemInput,
	...func(*dynamodb.Options),
) (*dynamodb.GetItemOutput, error) {
	return td.getOutput, td.errGet
}

func (td *testDynamoDB) PutItem(
	context.Context,
	*dynamodb.PutItemInput,
	...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return td.putOutput, td.errPut
}

var (
	errTest = errors.New("error")
	testID  = "test-id"
)

func fnID(_ string) string {
	return testID
}

func TestErrors(t *testing.T) {

	repo := New(&testDynamoDB{
		errDelete: errTest,
		errGet:    errTest,
		errPut:    errTest,
	},
		fnID,
	)
	res, err := repo.GetUser(context.Background(), testID)

	assert.Nil(t, res)
	assert.Error(t, err)

	err = repo.PutUser(context.Background(), types.User{})

	assert.Error(t, err)

	err = repo.RemoveUser(context.Background(), testID)

	assert.Error(t, err)
}
