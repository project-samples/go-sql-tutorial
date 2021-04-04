package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/common-go/config"
	"github.com/common-go/log"
	m "github.com/common-go/middleware"
	"github.com/gorilla/mux"

	"go-service/internal/app"
)

func main() {
	var conf app.Root
	er1 := config.Load(&conf, "configs/config")
	if er1 != nil {
		panic(er1)
	}

	r := mux.NewRouter()

	log.Initialize(conf.Log)
	r.Use(m.BuildContext)
	logger := m.NewStructuredLogger()
	r.Use(m.Logger(conf.MiddleWare, log.InfoFields, logger))
	r.Use(m.Recover(log.ErrorMsg))

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
