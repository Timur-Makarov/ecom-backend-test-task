package app

import (
	bannerRepos "ecom-backend-test-task/internal/banner/repositories"
	"github.com/jackc/pgx/v5"
	"gorm.io/gorm"
)

type Repositories struct {
	PGGBannerRepository *bannerRepos.PGGBannerRepository
	PGCBannerRepository *bannerRepos.PGCBannerRepository
}

func GetPGGRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		PGGBannerRepository: &bannerRepos.PGGBannerRepository{DB: db},
	}
}

func GetPGCRepositories(db *pgx.Conn) *Repositories {
	return &Repositories{
		PGCBannerRepository: &bannerRepos.PGCBannerRepository{DB: db},
	}
}
