package repositories

import (
	"context"
	"ecom-backend-test-task/internal/banner/domain"
	"ecom-backend-test-task/internal/pkg/database/pgg"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type PGGBannerRepository struct {
	DB *gorm.DB
}

func (r PGGBannerRepository) CreateBanner(name string) error {
	pggBanner := pgg.Banner{
		Name: name,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.DB.WithContext(ctx).Create(&pggBanner).Error
}

func (r PGGBannerRepository) CreateOrUpdateCounterStatistics(stats map[int]map[int32]domain.CounterStatistic) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var allCounterStats []pgg.CounterStatistic

	for _, value := range stats {
		for _, v := range value {
			ppgCS := pgg.CounterStatistic{
				BannerID:      v.BannerID,
				TimestampTo:   v.TimestampTo,
				TimestampFrom: v.TimestampFrom,
				Count:         v.Count,
			}
			allCounterStats = append(allCounterStats, ppgCS)
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

func (r PGGBannerRepository) GetBannerCounterStatistics(bannerID int32, tsFrom, tsTo int64) ([]domain.CounterStatistic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var stats []pgg.CounterStatistic

	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var banner pgg.Banner

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

	var convertedStats []domain.CounterStatistic

	for _, stat := range stats {
		convertedStats = append(convertedStats, domain.CounterStatistic{
			BannerID:      stat.BannerID,
			TimestampTo:   stat.TimestampTo,
			TimestampFrom: stat.TimestampFrom,
			Count:         stat.Count,
		})
	}

	return convertedStats, nil
}
