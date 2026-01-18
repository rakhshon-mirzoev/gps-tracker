package main

import (
	"fmt"
	"log"

	"github.com/rakhshon-mirzoev/gps-tracker/config"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/app"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		return err
	}

	app, err := app.New(cfg)
	if err != nil {
		return err
	}

	return app.Run(fmt.Sprintf("%s:%s", cfg.Core.Host, cfg.Core.Port))
}
