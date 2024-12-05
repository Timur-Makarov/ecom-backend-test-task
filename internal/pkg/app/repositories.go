package app

import (
	bannerRepos "ecom-backend-test-task/internal/banner/repositories"
	"gorm.io/gorm"
)

type Repositories struct {
	BannerRepository *bannerRepos.BannerRepository
}

func GetRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		BannerRepository: &bannerRepos.BannerRepository{DB: db},
	}
}
