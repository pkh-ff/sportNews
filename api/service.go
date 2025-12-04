package api

import (
	"context"
	"fmt"
	"net/http"
	"sportNews/conf"
	"sportNews/internal/service"
	"sportNews/pkg/log"
	"time"

	"go.uber.org/zap"
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

	e := engine(conf.Debug)
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
	log.Info("Starting HTTP server", zap.String("address", serv.HttpServer.Addr))
	if err := serv.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("Failed to start HTTP server", zap.Error(err))
	} else {
		log.Info("HTTP server stopped")
	}
}

func (serv *Server) Shutdown(ctx context.Context) error {
	log.Info("Shutting down HTTP server")
	err := serv.HttpServer.Shutdown(ctx)
	if err != nil {
		log.Error("HTTP server shutdown error", zap.Error(err))
		return err
	}

	log.Info("HTTP server shutdown successfully")
	return nil
}
