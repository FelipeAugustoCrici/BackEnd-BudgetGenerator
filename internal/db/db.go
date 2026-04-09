package db

import (
	"log"
	"os"

	"budgetgen/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	err = DB.AutoMigrate(
		&model.User{},
		&model.Quote{},
		&model.Template{},
		&model.Settings{},
	)
	if err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	log.Println("database connected and migrated")
}
