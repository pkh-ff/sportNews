package api

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"sportNews/conf"
	"sportNews/internal/log"
	"sportNews/internal/service"
	"time"
	"xorm.io/xorm"
)

type App struct {
	Ctx  context.Context
	Serv *service.Serv
}

type Server struct {
	HttpServer *http.Server
}

func New(ctx context.Context, conf *conf.App, db *xorm.EngineGroup) *Server {
	app := App{
		Ctx:  ctx,
		Serv: service.New(db),
	}

	e := engine()
	app.router(e)
	app.Serv.Repo.NewDBSession()
	addr := fmt.Sprintf(":%s", conf.Addr)
	h := &http.Server{
		Addr:         addr,
		Handler:      e,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return &Server{
		HttpServer: h,
	}
}

func (serv *Server) Run() {
	if err := serv.HttpServer.ListenAndServe(); err != nil {
		log.Error("run http server", zap.Error(err))
	}
}

func (serv *Server) Shutdown(ctx context.Context) error {
	return serv.Shutdown(ctx)
}
