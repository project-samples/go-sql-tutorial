package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"go-service/internal/app"
	"go-service/pkg/config"
)

func main() {
	var conf app.Root
	er1 := config.Load(&conf, "configs", "config")
	if er1 != nil {
		panic(er1)
	}

	r := mux.NewRouter()

	er2 := app.Route(r, context.Background(), conf.DB)
	if er2 != nil {
		panic(er2)
	}
	fmt.Println("Start server")
	server := ""
	if conf.Server.Port > 0 {
		server = ":" + strconv.Itoa(conf.Server.Port)
	}
	http.ListenAndServe(server, r)
}
