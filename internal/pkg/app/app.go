package app

import (
	"ecom-backend-test-task/internal/pkg/database"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
)

type App struct {
	repositories *Repositories
	handlers     *Handlers
	services     *Services
	closeCh      chan os.Signal
	fatalCh      chan error
	closers      []func() error
	server       *http.Server
	logger       *slog.Logger
}

func NewApp() (*App, error) {
	hardCodedDSN := "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"

	var err error

	app := &App{
		fatalCh: make(chan error),
	}

	db, err := app.initDB(hardCodedDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to init db connection: %w", err)
	}

	shouldMigrate := app.checkIfShouldMigrate()

	if shouldMigrate {
		err := database.MigrateDB(db)
		if err != nil {
			return nil, fmt.Errorf("failed to run db migrations: %w", err)
		}
	}

	err = app.initServer(db)
	if err != nil {
		return nil, fmt.Errorf("failed to init http server: %w", err)
	}

	err = app.initGracefulShutdown()
	if err != nil {
		return nil, fmt.Errorf("failed to init graceful shutdown: %w", err)
	}

	return app, nil
}

func (a *App) Run() error {
	go func() {
		a.logger.Info("Server is running on :8080")
		err := a.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("http server failed", "error", err)
		}
	}()

	select {
	case <-a.closeCh:
		a.logger.Info("App is shutting down")
		a.initClose()
	case err := <-a.fatalCh:
		a.logger.Error(err.Error())
		a.initClose()
	}

	return nil
}

func (a *App) initClose() {
	for i := len(a.closers) - 1; i >= 0; i-- {
		err := a.closers[i]()
		if err != nil {
			a.logger.Error("failed to close resource", "i", i, "error", err)
		}
	}
}

func (a *App) initLogger() {
	logger := slog.Default()
	a.logger = logger
}

func (a *App) initGracefulShutdown() error {
	a.closeCh = make(chan os.Signal, 1)
	signal.Notify(a.closeCh, os.Interrupt)

	return nil
}
