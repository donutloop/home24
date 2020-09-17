package api

import (
	"fmt"
	"github.com/donutloop/home24/internal/server"
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

func NewAPI(test bool) *API {
	return &API{
		Test: test,
	}
}

type API struct {
	addrs  string
	Server *server.Server
	Test   bool
}

func (a *API) Bootstrap() {
	a.Server = server.New()
	a.Server.InitHandlers()
}

func (a *API) Start() {

	logrus.Info("start server")
	if err := a.Server.Start(os.Getenv("SERVER_ADDRS"), a.Test); err != nil {
		log.Fatal(fmt.Sprintf("error server could not listen on addr %v, err: %v", a.addrs, err))
	}
}

func (a *API) Stop() {
	a.Server.Stop(a.Test)
}
