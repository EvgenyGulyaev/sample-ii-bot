package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/EvgenyGulyaev/sample-ii-bot/internal/app"
	"github.com/EvgenyGulyaev/sample-ii-bot/internal/config"
	"github.com/EvgenyGulyaev/sample-ii-bot/internal/envfile"
)

func main() {
	envfile.Load(".env")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	bot := app.New(cfg)
	if err := bot.Run(ctx); err != nil {
		log.Fatalf("bot stopped: %v", err)
	}
}
