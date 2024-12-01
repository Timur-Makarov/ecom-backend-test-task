package handlers

import (
	"ecom-backend-test-task/internal/services"
)

type Handlers struct {
	BannerHandler
}

func GetHandlers(services *services.Services) *Handlers {
	return &Handlers{
		BannerHandler{
			Services: services,
		},
	}
}
