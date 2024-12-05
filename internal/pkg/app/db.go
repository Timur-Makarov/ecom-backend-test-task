package app

import (
	"flag"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

func (a *App) initDB(dsn string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	for i := 1; i <= 5; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

		if err != nil && i == 5 {
			return nil, err
		}

		if err != nil {
			log.Printf("Could not connect to DB. Retry #%d \n", i)
			time.Sleep(1 * time.Second)
		}
	}

	return db, nil
}

func (a *App) checkIfShouldMigrate() bool {
	shouldMigrate := flag.Bool("runMigrations", false, "whether or not to run db migrations")
	flag.Parse()

	return *shouldMigrate
}
