package app

import (
	"context"
	"flag"
	"github.com/jackc/pgx/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

func (a *App) initGormDB(dsn string) (*gorm.DB, error) {
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

func (a *App) initSqlcDB(dsn string) (*pgx.Conn, error) {
	ctx, ctxClose := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctxClose()

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}

	a.closers = append(a.closers, func() error {
		return conn.Close(context.Background())
	})

	return conn, nil
}

func (a *App) checkIfShouldMigrate() bool {
	env := os.Getenv("ENVIRONMENT")
	if env != "development" {
		return false
	}

	shouldMigrate := flag.Bool("runMigrations", false, "whether or not to run db migrations")
	flag.Parse()

	return *shouldMigrate
}
