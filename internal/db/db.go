package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func New(databasePath string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("failed to get database instance")
	}

	_, err = sqlDB.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal("failed to enable foreign key constraints")
	}

	return db
}
