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

	return err
}

func (r PGCBannerRepository) CreateOrUpdateCounterStatistics(stats map[int]map[int32]domain.CounterStatistic) error {
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

	tx, err := r.DB.Begin(ctx)
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			slog.Error("failed to rollback transaction: ", err)
		}
	}(tx, ctx)
	if err != nil {
		return err
	}

	db := pgc.New(tx)
	res := db.CreateOrUpdateCounterStatistics(ctx, allCounterStats)

	res.Exec(func(_ int, e error) {
		if e != nil {
			err = e
		}
	})

	if err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r PGCBannerRepository) GetBannerCounterStatistics(bannerID int32, tsFrom, tsTo int64) ([]domain.CounterStatistic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := pgc.New(r.DB)

	_, err := db.GetBanner(ctx, bannerID)

	if err != nil {
		return nil, domain.GetNotFoundError("Banner not found")
	}

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
