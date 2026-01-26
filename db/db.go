package db

import (
	"log"
	"os"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	once     sync.Once
	instance *gorm.DB
)

func InitDB() {
	once.Do(func() {

		db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{
			Logger:                                   logger.Default.LogMode(logger.Info),
			DisableForeignKeyConstraintWhenMigrating: false,
		})
		if err != nil {
			log.Fatal("Error connecting to the database:", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatal("Error getting database connection:", err)
		}

		sqlDB.SetMaxOpenConns(25) // Set the maximum number of open connections
		sqlDB.SetMaxIdleConns(10) // Set the maximum number of idle connections

		log.Println("Connected to the database!")

		instance = db
	})
}

func GetDB() *gorm.DB {
	return instance
}
