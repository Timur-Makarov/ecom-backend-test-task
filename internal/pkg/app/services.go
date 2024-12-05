package app

import (
	bannerServices "ecom-backend-test-task/internal/banner/services"
)

type Services struct {
	BannerService *bannerServices.BannerService
}

func GetServices(repos *Repositories) *Services {
	return &Services{
		BannerService: &bannerServices.BannerService{
			Repo: repos.BannerRepository,
		},
	}
}
