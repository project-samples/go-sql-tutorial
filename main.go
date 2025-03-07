package main

import (
	"context"
	"fmt"
	"github.com/core-go/config"
	"github.com/core-go/core"
	"github.com/core-go/log"
	mid "github.com/core-go/log/middleware"
	"github.com/gorilla/mux"
	"net/http"

	"go-service/internal/app"
)

func main() {
	var conf app.Config
	er1 := config.Load(&conf, "configs/config")
	if er1 != nil {
		panic(er1)
	}

	r := mux.NewRouter()

	log.Initialize(conf.Log)
	r.Use(mid.BuildContext)
	logger := mid.NewLogger()
	if log.IsInfoEnable() {
		r.Use(mid.Logger(conf.MiddleWare, log.InfoFields, logger))
	}
	r.Use(mid.Recover(log.PanicMsg))

	er2 := app.Route(r, context.Background(), conf)
	if er2 != nil {
		panic(er2)
	}
	fmt.Println(core.ServerInfo(conf.Server))
	if er3 := http.ListenAndServe(core.Addr(conf.Server.Port), r); er3 != nil {
		fmt.Println(er3.Error())
	}
}
