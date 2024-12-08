package app

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func (a *App) initServer() error {
	if a.repositories == nil {
		return fmt.Errorf("app repositories should not be nil")
	}

	a.services = GetServices(a.repositories)
	a.Handlers = GetHandlers(a.services)

	http := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
	})

	http.Post("/banners", a.Handlers.BannerHandler.CreateBanner)
	http.Get("/banners/:bannerID/stats", a.Handlers.BannerHandler.GetCounterStatistics)
	http.Put("/banners/:bannerID/stats", a.Handlers.BannerHandler.UpdateCounterStatistics)

	a.Http = http

	go func() {
		if err := a.services.BannerService.RunCounterUpdater(); err != nil {
			a.logger.Error("RunCounterUpdater failed with error:", err)
			a.fatalCh <- err
		}
	}()

	a.closers = append(a.closers, func() error {
		return a.Http.Shutdown()
	})

	return nil
}

func errorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	err = ctx.JSON(fiber.Map{
		"error": fiber.Map{
			"code":    code,
			"message": err.Error(),
		},
	})

	if err != nil {
		ctx.Status(fiber.StatusInternalServerError)
	}

	return nil
}
