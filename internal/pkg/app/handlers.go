package app

import (
	"ecom-backend-test-task/internal/banner/handlers"
)

type Handlers struct {
	BannerHandler *handlers.BannerHandler
}

func GetHandlers(services *Services) *Handlers {
	return &Handlers{
		BannerHandler: &handlers.BannerHandler{
			Service: *services.BannerService,
		},
	}
}
