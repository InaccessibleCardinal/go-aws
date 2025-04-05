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
	db *users.UserRepo
}

func NewUsersController(db *users.UserRepo) *UsersController {
	return &UsersController{db: db}
}

func (u *UsersController) GetUser(w http.ResponseWriter, r *http.Request) {
	id := getID(r.URL.Path)
	user, err := u.db.GetUser(r.Context(), id)
	if err != nil {
		ErrorResponse(w, err)
		return
	}
	SuccessResponse(w, user)
}

func (u *UsersController) PutUser(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromBody(r.Body)
	if err != nil {
		ErrorResponse(w, err)
		return
	}
	err = u.db.PutUser(r.Context(), *user)
	if err != nil {
		ErrorResponse(w, err)
		return
	}
	CreatedResponse(w, map[string]string{"message": "user created"})
}

func getID(path string) string {
	parts := strings.Split(path, "/")
	return parts[1]
}

func getUserFromBody(body io.ReadCloser) (*types.User, error) {
	defer body.Close()
	bts, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	parsedUser, err := parseUser(bts)
	if err != nil {
		return nil, err
	}
	return parsedUser, nil
}

func parseUser(userBytes []byte) (*types.User, error) {
	var user types.User
	err := json.Unmarshal(userBytes, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
