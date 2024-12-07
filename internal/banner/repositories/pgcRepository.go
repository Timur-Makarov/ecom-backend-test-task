package repositories

import (
	"context"
	"ecom-backend-test-task/internal/banner/domain"
	"ecom-backend-test-task/internal/pkg/database/pgc"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"time"
)

type PGCBannerRepository struct {
	DB *pgx.Conn
}

func (r PGCBannerRepository) CreateBanner(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := pgc.New(r.DB)
	err := db.CreateBanner(ctx, name)

	if err != nil {
		slog.Error(err.Error())
	}

	return err
}

func (r PGCBannerRepository) CreateOrUpdateCounterStatistics(stats map[int]map[int32]domain.CounterStatistic) error {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var allCounterStats []pgc.CreateOrUpdateCounterStatisticsParams

	for _, value := range stats {
		for _, v := range value {
			params := pgc.CreateOrUpdateCounterStatisticsParams{
				BannerID:      v.BannerID,
				TimestampTo:   v.TimestampTo,
				TimestampFrom: v.TimestampFrom,
				Count:         v.Count,
			}
			allCounterStats = append(allCounterStats, params)
		}
	}

	db := pgc.New(r.DB)
	res := db.CreateOrUpdateCounterStatistics(ctx, allCounterStats)

	defer res.Close()

	res.Exec(func(_ int, e error) {
		if e != nil {
			err = e
		}
	})

	return err
}

func (r PGCBannerRepository) GetBannerCounterStatistics(bannerID int32, tsFrom, tsTo int64) ([]domain.CounterStatistic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := pgc.New(r.DB)
	stats, err := db.GetCounterStatistics(ctx, pgc.GetCounterStatisticsParams{
		BannerID:      bannerID,
		TimestampFrom: tsFrom,
		TimestampTo:   tsTo,
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
