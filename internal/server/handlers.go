package server

import (
	"encoding/json"
	"go-aws/internal/types"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	log.Printf("getting user by id: %s\n", id)
	user, err := s.userDBGet(r.Context(), "USER#"+id)
	if err != nil {
		ErrorResponse(w, err)
		return
	}
	SuccessResponse(w, user)
}

func (s *Server) PutUser(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromBody(r.Body)
	if err != nil {
		ErrorResponse(w, err)
		return
	}
	err = s.userDBPut(r.Context(), *user)
	if err != nil {
		ErrorResponse(w, err)
		return
	}
	CreatedResponse(w, map[string]string{"message": "user created"})
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
