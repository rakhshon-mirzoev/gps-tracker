package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/IkbolNazarov/gps-tracker/internal/websocket"
	"github.com/gin-gonic/gin"
)

type App struct {
	router   *gin.Engine
	server   *http.Server
	wsHub    *websocket.Hub
	wsCancel context.CancelFunc
}

func New(router *gin.Engine) *App {
	hub := websocket.NewHub()

	ctx, cancel := context.WithCancel(context.Background())
	go hub.Run(ctx)
	return &App{
		router:   router,
		wsHub:    hub,
		wsCancel: cancel,
	}
}

func (a *App) Run(addr string) error {
	a.server = &http.Server{
		Addr:        addr,
		Handler:     a.router,
		ReadTimeout: 10 * time.Second,
		IdleTimeout: 120 * time.Second,
	}

	log.Printf("server started on %s", addr)
	return a.server.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	if a.wsCancel != nil {
		a.wsCancel()
	}

	if a.server == nil {
		return errors.New("server is nil")
	}
	return a.server.Shutdown(ctx)
}

func (a *App) WsHub() *websocket.Hub {
	return a.wsHub
}
