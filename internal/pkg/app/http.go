package app

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func (a *App) initServer() error {
	if a.repositories == nil {
		return fmt.Errorf("app repositories should not be nil")
	}

	a.services = GetServices(a.repositories)
	a.Handlers = GetHandlers(a.services)

	http := fiber.New()

	http.Post("/banners", a.Handlers.BannerHandler.CreateBanner)
	http.Put("/banners/:bannerID/stats", a.Handlers.BannerHandler.UpdateCounterStatistics)
	http.Get("/banners/:bannerID/stats", a.Handlers.BannerHandler.GetCounterStatistics)

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
