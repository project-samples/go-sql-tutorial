package app

import (
	"context"
	"github.com/gorilla/mux"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	PATCH  = "PATCH"
	DELETE = "DELETE"
)

func Route(r *mux.Router, ctx context.Context, root Config) error {
	app, err := NewApp(ctx, root)
	if err != nil {
		return err
	}

	r.HandleFunc("/health", app.Health.Check).Methods(GET)

	userPath := "/users"
	r.HandleFunc(userPath, app.User.All).Methods(GET)
	r.HandleFunc(userPath+"/{id}", app.User.Load).Methods(GET)
	r.HandleFunc(userPath, app.User.Insert).Methods(POST)
	r.HandleFunc(userPath+"/{id}", app.User.Update).Methods(PUT)
	r.HandleFunc(userPath+"/{id}", app.User.Patch).Methods(PATCH)
	r.HandleFunc(userPath+"/{id}", app.User.Delete).Methods(DELETE)

	return nil
}
