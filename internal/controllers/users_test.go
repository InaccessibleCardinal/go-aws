package controllers

import (
	"errors"
	"go-aws/internal/parsers"
	"go-aws/internal/types"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
)

var (
	goodId = "goodId"
	badId = "badId"
)
type MockDb struct {
	mock.Mock
}

func (m *MockDb) GetUser(userId string) (*types.UserRecord, error) {
	args := m.Called(userId)
	return args.Get(0).(*types.UserRecord), args.Error(1)
}

func (m *MockDb) PutUser(user types.UserRecord)  error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockDb) RemoveUser(userId string) error {
	return nil
}


func Test_GetUserSuccess(t *testing.T) {
	testUser := types.UserRecord{UserId: goodId, UserName: "testKen", Email: "me@site.com"}
	testDb := new(MockDb)
	testDb.On("GetUser", goodId).Return(&testUser, nil)

	req := httptest.NewRequest(http.MethodGet, "/" + goodId, nil)
	res := httptest.NewRecorder()

	usersController := NewUsersController(testDb)
	usersController.GetUser(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected 200 but got %d", res.Code)
	}
	
	parsedUserFromBody := parsers.MustParseJson[types.UserRecord](res.Body.Bytes())

	if parsedUserFromBody.Email != testUser.Email {
		t.Fatalf("expected %s but got %s", testUser.Email, parsedUserFromBody.Email)
	}

	if parsedUserFromBody.UserName != testUser.UserName {
		t.Fatalf("expected %s but got %s", testUser.UserName, parsedUserFromBody.UserName)
	}
}

func Test_GetUserFailure(t *testing.T) {
	testEmptyUser := types.UserRecord{}
	testDb := new(MockDb)
	testDb.On("GetUser", badId).Return(&testEmptyUser, errors.New("do you even go?"))

	req := httptest.NewRequest(http.MethodGet, "/" + badId, nil)
	res := httptest.NewRecorder()

	usersController := NewUsersController(testDb)
	usersController.GetUser(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 but got %d", res.Code)
	}
}

func Test_PutUserSuccess(t *testing.T) {
	testNewUserBody := `{"email":"me@mail.com","userName":"testKen"}`
	testNewUser := parsers.MustParseJson[types.UserRecord]([]byte(testNewUserBody))
	testDb := new(MockDb)
	testDb.On("PutUser", testNewUser).Return(nil)

	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(testNewUserBody))
	res := httptest.NewRecorder()

	usersController := NewUsersController(testDb)
	usersController.PutUser(res, req)

	if res.Code != http.StatusCreated {
		t.Fatalf("expected 201 but got %d", res.Code)
	}
}