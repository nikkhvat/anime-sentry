package main

import (
	"anime-sentry/handlers/bot"
	"anime-sentry/pkg/scheduler"
	"context"
	"os"
	"os/signal"
	"syscall"

	"anime-sentry/repositories/postgres"

	"log"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

func main() {
	db, _ := postgres.New()

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	ctx := context.TODO()

	// * Launch a goroutine for regular status checks (every 10 minutes)
	scheduler.CheckNewEpisodes(ctx, db)

	// * Launch bot
	bot := bot.New(db)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		err := bot.Start(ctx)
		if err != nil {
			log.Printf("Bot failed to start: %v", err)
			os.Exit(1)
		}
	}()

	<-sigChan

	log.Println("Received an interrupt, Bot stopped...")
}
