package app

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type App struct {
	router *gin.Engine
	server *http.Server
}

func New(router *gin.Engine) *App {
	return &App{router: router}
}

func (a *App) Run(addr string) error {
	a.server = &http.Server{
		Addr:    addr,
		Handler: a.router,
	}

	log.Printf("server started on %s", addr)
	return a.server.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	if a.server == nil {
		return errors.New("server is nil.")
	}
	return a.server.Shutdown(ctx)
}
