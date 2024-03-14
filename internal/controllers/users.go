package controllers

import (
	"encoding/json"
	"go-aws/internal/repos/users"
	"go-aws/internal/types"
	"io"
	"net/http"
	"strings"
)

type UsersController struct {
	db users.UserDbIface
}

func NewUsersController(db users.UserDbIface) *UsersController {
	return &UsersController{db: db}
}

func (u *UsersController) GetUser(w http.ResponseWriter, r *http.Request) {
	id := getId(r.URL.Path)
	user, err := u.db.GetUser(id)
	if err != nil {
		ErrorResponse(w, err)
		return
	}
	SuccessResponse(w, user)
}

func (u *UsersController) PutUser(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromBody(r.Body, io.ReadAll)
	if err != nil {
		ErrorResponse(w, err)
		return
	}
	err = u.db.PutUser(*user)
	if err != nil {
		ErrorResponse(w, err)
		return
	}
	CreatedResponse(w, map[string]string{"message": "user created"})
}

func getId(path string) string {
	parts := strings.Split(path, "/")
	return parts[1]
}

func getUserFromBody(body io.ReadCloser, reader BodyReader) (*types.UserRecord, error) {
	defer body.Close()
	bts, err := reader(body)
	if err != nil {
		return nil, err
	}
	parsedUser, err := parseUser(bts)
	if err != nil {
		return nil, err
	}
	return parsedUser, nil
}

func parseUser(userBytes []byte) (*types.UserRecord, error) {
	var user types.UserRecord
	err := json.Unmarshal(userBytes, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
