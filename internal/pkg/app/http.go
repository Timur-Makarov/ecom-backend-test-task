package app

import (
	"context"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func (a *App) initServer(db *gorm.DB) error {
	a.repositories = GetRepositories(db)
	a.services = GetServices(a.repositories)
	a.Handlers = GetHandlers(a.services)

	server := &http.Server{
		Addr: "127.0.0.1:8080",
	}

	http.HandleFunc("/banner/add", a.Handlers.BannerHandler.AddBanner)
	http.HandleFunc("/counter/{bannerID}", a.Handlers.BannerHandler.UpdateCounterStats)
	http.HandleFunc("/stats/{bannerID}", a.Handlers.BannerHandler.GetCounterStats)

	a.server = server

	a.services.BannerService.RunCounterUpdater(a.fatalCh)

	a.closers = append(a.closers, func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		return a.server.Shutdown(ctx)
	})

	return nil
}
