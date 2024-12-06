package app

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func (a *App) initServer(db *gorm.DB) error {
	a.repositories = GetRepositories(db)
	a.services = GetServices(a.repositories)
	a.Handlers = GetHandlers(a.services)

	http := fiber.New()

	http.Post("/banners", a.Handlers.BannerHandler.AddBanner)
	http.Put("/banners/:bannerID/stats", a.Handlers.BannerHandler.UpdateCounterStats)
	http.Get("/banners/:bannerID/stats", a.Handlers.BannerHandler.GetCounterStats)

	a.Http = http

	go func() {
		if err := a.services.BannerService.RunCounterUpdater(); err != nil {
			a.fatalCh <- err
		}
	}()

	a.closers = append(a.closers, func() error {
		return a.Http.Shutdown()
	})

	return nil
}
