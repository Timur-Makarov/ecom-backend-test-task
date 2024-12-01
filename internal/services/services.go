package services

import "ecom-backend-test-task/internal/database"

type Services struct {
	BannerService
}

func GetServices(dbRepo *database.DBRepository) *Services {
	return &Services{
		BannerService{
			DBRepo:             dbRepo,
			statsByMinuteCache: make(map[int]map[uint64]database.CounterStats),
		},
	}
}
