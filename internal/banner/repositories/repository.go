package repositories

import (
	"context"
	"ecom-backend-test-task/internal/pkg/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type BannerRepository struct {
	DB *gorm.DB
}

func (r BannerRepository) SaveBanner(newBanner database.Banner) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.DB.WithContext(ctx).Create(&newBanner).Error
}

func (r BannerRepository) UpdateOrCreateBannerCounterStats(stats map[int]map[uint64]database.CounterStats) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var allCounterStats []database.CounterStats

	for _, value := range stats {
		for _, v := range value {
			allCounterStats = append(allCounterStats, v)
		}
	}

	query := r.DB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "timestamp_from"},
				{Name: "timestamp_to"},
			},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"count": gorm.Expr("counter_stats.count + EXCLUDED.count"),
			}),
		}).
		Create(&allCounterStats)

	return query.Error
}

func (r BannerRepository) GetBannerCounterStats(bannerID, tsFrom, tsTo uint64) ([]database.CounterStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var stats []database.CounterStats

	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var banner database.Banner

		if err := tx.First(&banner, bannerID).Error; err != nil {
			return err
		}

		err := tx.
			Where("timestamp_from >= ?", tsFrom).
			Where("timestamp_to <= ?", tsTo).
			Where("banner_id = ?", bannerID).
			Find(&stats).Error

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return stats, nil
}
