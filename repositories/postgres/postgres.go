package postgres

import (
	env "anime-sentry/pkg/env"
	repositories "anime-sentry/repositories"

	postgresDriver "gorm.io/driver/postgres"

	"fmt"
	"log"

	"anime-sentry/models"

	"gorm.io/gorm"
)

type postgres struct {
	db *gorm.DB
}

func New() (repositories.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		env.Get("POSTGRES_HOST"),
		env.Get("POSTGRES_USER"),
		env.Get("POSTGRES_PASSWORD"),
		env.Get("POSTGRES_DB_NAME"),
		env.Get("POSTGRES_PORT"),
		env.Get("POSTGRES_SSL_MODE"),
		env.Get("POSTGRES_TIMEZONE"),
	)

	db, err := gorm.Open(postgresDriver.New(postgresDriver.Config{DSN: dsn}), &gorm.Config{})

	if err != nil {
		log.Panicln(err)
	}

	_ = db.AutoMigrate(&models.Anime{})
	_ = db.AutoMigrate(&models.Subscriber{})
	_ = db.AutoMigrate(&models.User{})

	return &postgres{db: db}, nil
}
