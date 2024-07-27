package scheduler

import (
	"anime-sentry/repositories"
	"anime-sentry/services/anime"
	"context"
	"log"
	"time"

	"github.com/go-co-op/gocron"
)

func CheckNewEpisodes(ctx context.Context, db repositories.DB) {
	s := gocron.NewScheduler(time.UTC)

	anime := anime.New(db)

	_, err := s.Every(10).Minute().Do(func() {
		anime.CheckAnimeStatus(ctx)
	})

	if err != nil {
		log.Fatal(err)
	}
}
