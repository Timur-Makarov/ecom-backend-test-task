package config

import (
	"ecom-backend-test-task/internal/database"
	"ecom-backend-test-task/internal/handlers"
	"ecom-backend-test-task/internal/services"
	"flag"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

type APIComponents struct {
	Handlers *handlers.Handlers
	Services *services.Services
}

func SetupAPIComponents(db *gorm.DB) *APIComponents {
	dbRepo := database.GetDBRepository(db)
	s := services.GetServices(dbRepo)
	h := handlers.GetHandlers(s)

	return &APIComponents{
		h,
		s,
	}
}

func OpenDB(dsn string) (*gorm.DB, error) {
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

func CheckIfShouldMigrate(db *gorm.DB) {
	shouldMigrate := flag.Bool("runMigrations", false, "whether or not to run db migrations")
	flag.Parse()

	if *shouldMigrate {
		err := database.MigrateDB(db)
		if err != nil {
			log.Fatalf("could not run db migrations. Error - %v", err.Error())
		}
		log.Println("Successfully migrated DB")
		return
	}
}
