package server

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

const (
	applicationJSON = "application/json"
)

func TestServer_200s(t *testing.T) {
	testDB := &testDynamoDB{getOutput: &dynamodb.GetItemOutput{Item: map[string]dynamodbTypes.AttributeValue{
		"USERID": &dynamodbTypes.AttributeValueMemberS{Value: "tid"},
	}}}
	router := chi.NewMux()
	app := Server{
		db:          testDB,
		idGenerator: fnID,
		router:      router,
		server:      &http.Server{Addr: addr, Handler: router},
		userTable:   "users",
	}

	type testCase struct {
		name       string
		body       io.Reader
		requestURL string
		codeWant   int
	}

	go app.Run()

	for _, tc := range []testCase{
		{
			name:       "home",
			requestURL: "http://localhost" + addr + "/",
			codeWant:   http.StatusOK,
		},
		{
			name:       "users/tid",
			requestURL: "http://localhost" + addr + "/users/tid",
			codeWant:   http.StatusOK,
		},
		{
			name:       "post user error",
			body:       strings.NewReader(`{"userName":"test-user","email":"test@test.com"}`),
			requestURL: "http://localhost" + addr + "/users",
			codeWant:   http.StatusCreated,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var res *http.Response
			var err error
			if tc.body != nil {
				res, err = http.Post(tc.requestURL, applicationJSON, tc.body)
			} else {
				res, err = http.Get(tc.requestURL)
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.codeWant, res.StatusCode)
		})
	}

	t.Cleanup(func() {
		if err := app.server.Shutdown(context.Background()); err != nil {
			t.Fatalf("error shutting down server %s", err)
		}
	})
}

func TestServer_400(t *testing.T) {
	testDB := &testDynamoDB{errGet: errors.New("error getting"), errPut: errors.New("error putting")}
	router := chi.NewMux()
	app := Server{
		db:          testDB,
		idGenerator: fnID,
		router:      router,
		server:      &http.Server{Addr: addr, Handler: router},
		userTable:   "users",
	}

	type testCase struct {
		name       string
		body       io.Reader
		requestURL string
		codeWant   int
	}

	go app.Run()

	for _, tc := range []testCase{
		{
			name:       "post user error",
			body:       strings.NewReader(`{"userName":"test-user","email":"test@test.com"}`),
			requestURL: "http://localhost" + addr + "/users",
			codeWant:   http.StatusBadRequest,
		},
		{
			name:       "users/tid error",
			requestURL: "http://localhost" + addr + "/users/tid",
			codeWant:   http.StatusBadRequest,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var res *http.Response
			var err error
			if tc.body != nil {
				res, err = http.Post(tc.requestURL, applicationJSON, tc.body)
			} else {
				res, err = http.Get(tc.requestURL)
			}
			assert.Nil(t, err)
			assert.Equal(t, tc.codeWant, res.StatusCode)
		})
	}

	t.Cleanup(func() {
		if err := app.server.Shutdown(context.Background()); err != nil {
			t.Fatalf("error shutting down server %s", err)
		}
	})
}

func Test_putNilUser(t *testing.T) {
	testDB := &testDynamoDB{}
	router := chi.NewMux()
	app := Server{
		db:          testDB,
		idGenerator: fnID,
		router:      router,
		server:      &http.Server{Addr: addr, Handler: router},
		userTable:   "users",
	}

	go app.Run()

	res, err := http.Post("http://localhost"+addr+"/users", applicationJSON, nil)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	t.Cleanup(func() {
		if err := app.server.Shutdown(context.Background()); err != nil {
			t.Fatalf("error shutting down server %s", err)
		}
	})
}

type badBody string

func (b badBody) Read(p []byte) (n int, err error) {
	return 0, errors.New(string(b))
}

func (b badBody) Close() error {
	return nil
}

func Test_getUserFromBody(t *testing.T) {
	var body badBody = "error reading"
	user, err := getUserFromBody(body)

	assert.Nil(t, user)
	assert.EqualError(t, err, "error reading")
}

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
