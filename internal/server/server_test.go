package server

import (
	"context"
	"go-aws/internal/controllers"
	"go-aws/internal/repos/users"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {

	app := New(controllers.NewUsersController(users.InitUsersDB(context.Background())))

	go app.Run()

	res, err := http.Get("http://localhost:8080")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	t.Cleanup(func() {
		if err := app.server.Shutdown(context.Background()); err != nil {
			t.Fatalf("error shutting down server %s", err)
		}
	})
}
