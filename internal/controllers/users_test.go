package controllers

import (
	"context"
	"errors"
	"go-aws/internal/parsers"
	"go-aws/internal/repos/users"
	"go-aws/internal/types"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	badID      = "badId"
	goodID     = "goodId"
	getItem    = "GetItem"
	deleteItem = "DeleteItem"
	putItem    = "PutItem"
)

func Test_GetUser_Success(t *testing.T) {
	testUser := types.User{UserID: "id1", UserName: "test user", Email: "test@mail.com"}
	testOutput := goodUserFromDB(testUser)
	testDB := new(MockUsersDB)
	testDB.On(getItem, mock.Anything, mock.Anything).Return(testOutput, nil).Once()
	usersRepo := users.New(testDB, testGenerator)

	req := httptest.NewRequest(http.MethodGet, "/id1", nil)
	res := httptest.NewRecorder()

	usersCtrl := NewUsersController(usersRepo)
	usersCtrl.GetUser(res, req)

	assert.Equal(t, http.StatusOK, res.Result().StatusCode)

	parsedUser := parsers.MustParseJson[types.User](res.Body.Bytes())

	assert.Equal(t, testUser, parsedUser)
}

func TestGetUser_DBError(t *testing.T) {
	var emptyUser *dynamodb.GetItemOutput
	testDB := new(MockUsersDB)
	testDB.On(getItem, mock.Anything, mock.Anything).Return(emptyUser, errors.New("no user")).Once()
	usersRepo := users.New(testDB, testGenerator)

	req := httptest.NewRequest(http.MethodGet, "/"+badID, nil)
	res := httptest.NewRecorder()

	usersController := NewUsersController(usersRepo)
	usersController.GetUser(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
}

func TestPutUser_Success(t *testing.T) {
	userBody := `{"email":"me@mail.com","userName":"testKen"}`
	testDB := new(MockUsersDB)
	testDB.On(putItem, mock.Anything, mock.Anything).Return(&dynamodb.PutItemOutput{}, nil).Once()
	usersRepo := users.New(testDB, testGenerator)

	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(userBody))
	res := httptest.NewRecorder()

	usersController := NewUsersController(usersRepo)
	usersController.PutUser(res, req)

	assert.Equal(t, http.StatusCreated, res.Code)
}

func TestPutUser_DBError(t *testing.T) {
	userBody := `{"email":"me@mail.com","userName":"testKen"}`
	testDB := new(MockUsersDB)
	testDB.On(putItem, mock.Anything, mock.Anything).Return(&dynamodb.PutItemOutput{}, errors.New("error")).Once()
	usersRepo := users.New(testDB, testGenerator)

	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(userBody))
	res := httptest.NewRecorder()

	usersController := NewUsersController(usersRepo)
	usersController.PutUser(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)
}

func TestPutUser_ClientError(t *testing.T) {
	userBody := `{"email":[1,2,3],"userName":"testKen"}`
	testDB := new(MockUsersDB)
	usersRepo := users.New(testDB, testGenerator)

	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(userBody))
	res := httptest.NewRecorder()

	usersController := NewUsersController(usersRepo)
	usersController.PutUser(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
}

type badReader struct {
	err error
}

func (br badReader) Read(p []byte) (n int, err error) { return 0, br.err }
func (br badReader) Close() error                     { return nil }

func Test_getUserFromBody_BadReader(t *testing.T) {
	user, err := getUserFromBody(badReader{err: errors.New("error")})

	assert.Nil(t, user)
	assert.Error(t, err)
}

func goodUserFromDB(user types.User) *dynamodb.GetItemOutput {
	userAV, _ := users.UserToAttributeValueMap(user)
	return &dynamodb.GetItemOutput{
		Item: userAV,
	}
}

func testGenerator(_ string) string {
	return "testing"
}

type MockUsersDB struct {
	mock.Mock
}

func (m *MockUsersDB) GetItem(
	context.Context,
	*dynamodb.GetItemInput,
	...func(*dynamodb.Options),
) (*dynamodb.GetItemOutput, error) {
	args := m.Called()
	return args.Get(0).(*dynamodb.GetItemOutput), args.Error(1)
}

func (m *MockUsersDB) PutItem(
	context.Context,
	*dynamodb.PutItemInput,
	...func(*dynamodb.Options),
) (*dynamodb.PutItemOutput, error) {
	args := m.Called()
	return args.Get(0).(*dynamodb.PutItemOutput), args.Error(1)
}

func (m *MockUsersDB) DeleteItem(
	context.Context,
	*dynamodb.DeleteItemInput,
	...func(*dynamodb.Options),
) (*dynamodb.DeleteItemOutput, error) {
	args := m.Called()
	return args.Get(0).(*dynamodb.DeleteItemOutput), args.Error(1)
}

func isContext(ctx context.Context) bool {
	return true
}
