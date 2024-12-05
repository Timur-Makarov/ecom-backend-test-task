package main

import (
	"ecom-backend-test-task/internal/pkg/app"
	"log/slog"
)

func main() {
	ap, err := app.NewApp()
	if err != nil {
		slog.Error("failed to create app", "error", err)
		return
	}

	if err = ap.Run(); err != nil {
		slog.Error("failed to run app", "error", err)
	}
}
