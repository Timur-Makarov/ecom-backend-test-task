package app

import (
	"ecom-backend-test-task/internal/pkg/database"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
)

type App struct {
	repositories *Repositories
	Handlers     *Handlers
	services     *Services
	closeCh      chan os.Signal
	// Channel for fatal errors in background goroutines
	fatalCh chan error
	closers []func() error
	Http    *fiber.App
	logger  *slog.Logger
}

func NewApp() (*App, error) {
	app := &App{
		logger:  slog.Default(),
		fatalCh: make(chan error),
	}

	err := app.initEnv()
	if err != nil {
		return nil, err
	}

	DSN := os.Getenv("DSN")
	if DSN == "" {
		return nil, fmt.Errorf("DSN environment variable not set")
	}

	db, err := app.initDB(DSN)
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
		return nil, fmt.Errorf("failed to init Http Http: %w", err)
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
		err := a.Http.Listen("127.0.0.1:8080")
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("Http Http failed", "error", err)
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

func (a *App) initGracefulShutdown() error {
	a.closeCh = make(chan os.Signal, 1)
	signal.Notify(a.closeCh, os.Interrupt)

	return nil
}

func (a *App) initEnv() error {
	envFilepath := "./.dev.env"

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %v", err)
	}

	if strings.HasSuffix(wd, "tests") {
		envFilepath = "../.test.env"
	}

	err = godotenv.Load(envFilepath)
	if err != nil {
		return fmt.Errorf("failed to load %s file: %v", envFilepath, err)
	}

	return nil
}
