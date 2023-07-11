package database

import (
	"fmt"

	"log"

	models "anime-bot-schedule/models"
	env "anime-bot-schedule/pkg/env"

	postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		env.Get("POSTGRES_HOST"),
		env.Get("POSTGRES_USER"),
		env.Get("POSTGRES_PASSWORD"),
		env.Get("POSTGRES_DB_NAME"),
		env.Get("POSTGRES_PORT"),
		env.Get("POSTGRES_SSL_MODE"),
		env.Get("POSTGRES_TIMEZONE"),
	)

	db, err := gorm.Open(postgres.New(postgres.Config{DSN: dsn}), &gorm.Config{})

	if err != nil {
		log.Panicln(err)
	}

	autoMigrateDB(db)

	return db
}

func autoMigrateDB(db *gorm.DB) {
	_ = db.AutoMigrate(&models.Anime{})
	_ = db.AutoMigrate(&models.Subscriber{})
}
