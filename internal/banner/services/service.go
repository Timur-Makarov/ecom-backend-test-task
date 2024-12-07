package services

import (
	"ecom-backend-test-task/internal/banner/domain"
	"fmt"
	"sync"
	"time"
)

var mutex = &sync.Mutex{}
var statsByMinuteCache = make(map[int]map[int32]domain.CounterStatistic)

type BannerRepository interface {
	CreateBanner(name string) error
	CreateOrUpdateCounterStatistics(stats map[int]map[int32]domain.CounterStatistic) error
	GetBannerCounterStatistics(bannerID int32, tsFrom int64, tsTo int64) ([]domain.CounterStatistic, error)
}

type BannerService struct {
	Repo BannerRepository
}

func (s *BannerService) RunCounterUpdater() error {
	for {
		time.Sleep(10 * time.Second)
		if len(statsByMinuteCache) != 0 {
			err := s.Repo.CreateOrUpdateCounterStatistics(statsByMinuteCache)
			if err != nil {
				return fmt.Errorf("failed to save stats into db: %v \n", err)
			}
			mutex.Lock()
			statsByMinuteCache = make(map[int]map[int32]domain.CounterStatistic)
			mutex.Unlock()
		}
	}
}

func (s *BannerService) CreateBanner(name string) error {
	return s.Repo.CreateBanner(name)
}

func (s *BannerService) UpdateBannerCounterStats(bannerID int32) {
	ts := time.Now()
	minute := ts.Minute()
	timestampFrom := ts.Truncate(time.Minute).Unix()
	timestampTo := ts.Truncate(time.Minute).Add(time.Minute).Unix() - 1

	mutex.Lock()
	if statsByMinuteCache[minute] == nil {
		statsByMinuteCache[minute] = make(map[int32]domain.CounterStatistic)
	}

	if statsByMinuteCache[minute][bannerID].Count == 0 {
		newCounterStats := domain.CounterStatistic{
			BannerID:      bannerID,
			TimestampFrom: timestampFrom,
			TimestampTo:   timestampTo,
			Count:         1,
		}

		statsByMinuteCache[minute][bannerID] = newCounterStats
	} else {
		counterStats := statsByMinuteCache[minute][bannerID]
		counterStats.Count += 1
		statsByMinuteCache[minute][bannerID] = counterStats
	}
	mutex.Unlock()
}

type GetCounterStatsDTO struct {
	BannerID      int32 `json:"bannerId"`
	TimestampFrom int64 `json:"timestampFrom"`
	TimestampTo   int64 `json:"timestampTo"`
	Count         int64 `json:"count"`
}

func (s *BannerService) GetCounterStats(bannerID int32, tsFrom, tsTo int64) (*GetCounterStatsDTO, error) {
	counterStats, err := s.Repo.GetBannerCounterStatistics(bannerID, tsFrom, tsTo)

	if err != nil {
		return nil, err
	}

	var count int64

	for _, stat := range counterStats {
		count += stat.Count
	}

	response := &GetCounterStatsDTO{
		BannerID:      bannerID,
		TimestampFrom: tsFrom,
		TimestampTo:   tsTo,
		Count:         count,
	}

	return response, nil
}
