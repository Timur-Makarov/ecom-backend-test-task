package services

import (
	"ecom-backend-test-task/internal/pkg/database"
	"fmt"
	"sync"
	"time"
)

var mutex = &sync.Mutex{}
var statsByMinuteCache = make(map[int]map[uint64]database.CounterStats)

type BannerRepository interface {
	SaveBanner(newBanner database.Banner) error
	UpdateOrCreateBannerCounterStats(stats map[int]map[uint64]database.CounterStats) error
	GetBannerCounterStats(bannerID uint64, tsFrom uint64, tsTo uint64) ([]database.CounterStats, error)
}

type BannerService struct {
	Repo BannerRepository
}

func (s *BannerService) RunCounterUpdater() error {
	for {
		time.Sleep(10 * time.Second)
		if len(statsByMinuteCache) != 0 {
			err := s.Repo.UpdateOrCreateBannerCounterStats(statsByMinuteCache)
			if err != nil {
				return fmt.Errorf("failed to save stats into db: %v \n", err)
			}
			mutex.Lock()
			statsByMinuteCache = make(map[int]map[uint64]database.CounterStats)
			mutex.Unlock()
		}
	}
}

func (s *BannerService) AddBanner(name string) error {
	newBanner := database.Banner{
		Name: name,
	}

	return s.Repo.SaveBanner(newBanner)
}

func (s *BannerService) UpdateBannerCounterStats(bannerID uint64) {
	ts := time.Now()
	minute := ts.Minute()
	timestampFrom := ts.Truncate(time.Minute).Unix()
	timestampTo := ts.Truncate(time.Minute).Add(time.Minute).Unix() - 1

	mutex.Lock()
	if statsByMinuteCache[minute] == nil {
		statsByMinuteCache[minute] = make(map[uint64]database.CounterStats)
	}

	if statsByMinuteCache[minute][bannerID].Count == 0 {
		newCounterStats := database.CounterStats{
			BannerID:      bannerID,
			TimestampFrom: uint64(timestampFrom),
			TimestampTo:   uint64(timestampTo),
			Count:         uint64(1),
		}

		statsByMinuteCache[minute][bannerID] = newCounterStats
	} else {
		counterStats := statsByMinuteCache[minute][bannerID]
		counterStats.Count += uint64(1)
		statsByMinuteCache[minute][bannerID] = counterStats
	}
	mutex.Unlock()
}

type GetCounterStatsDTO struct {
	BannerID      uint64 `json:"bannerId"`
	Count         uint64 `json:"count"`
	TimestampFrom uint64 `json:"timestampFrom"`
	TimestampTo   uint64 `json:"timestampTo"`
}

func (s *BannerService) GetCounterStats(bannerID uint64, tsFrom, tsTo uint64) (*GetCounterStatsDTO, error) {
	counterStats, err := s.Repo.GetBannerCounterStats(bannerID, tsFrom, tsTo)

	if err != nil {
		return nil, err
	}

	var count uint64

	for _, stat := range counterStats {
		count += stat.Count
	}

	response := &GetCounterStatsDTO{
		BannerID:      bannerID,
		Count:         count,
		TimestampFrom: tsFrom,
		TimestampTo:   tsTo,
	}

	return response, nil
}
