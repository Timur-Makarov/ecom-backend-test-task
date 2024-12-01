package services

import (
	"ecom-backend-test-task/internal/database"
	"log"
	"sync"
	"time"
)

var mutex = &sync.Mutex{}

type BannerService struct {
	DBRepo             *database.DBRepository
	statsByMinuteCache map[int]map[uint64]database.CounterStats
}

func (s *BannerService) RunCounterUpdater() {
	go func() {
		for {
			time.Sleep(10 * time.Second)
			if len(s.statsByMinuteCache) != 0 {
				err := s.DBRepo.UpdateOrCreateBannerCounterStats(s.statsByMinuteCache)
				if err != nil {
					log.Fatalf("error while saving stats into db: %v \n", err)
				}
				mutex.Lock()
				s.statsByMinuteCache = make(map[int]map[uint64]database.CounterStats)
				mutex.Unlock()
			}
		}
	}()
}

func (s *BannerService) AddBanner(name string) error {
	newBanner := database.Banner{
		Name: name,
	}

	return s.DBRepo.SaveBanner(newBanner)
}

func (s *BannerService) UpdateBannerCounterStats(bannerID uint64) {
	ts := time.Now()
	minute := ts.Minute()
	timestampFrom := ts.Truncate(time.Minute).Unix()
	timestampTo := ts.Truncate(time.Minute).Add(time.Minute).Unix() - 1

	mutex.Lock()
	if s.statsByMinuteCache[minute] == nil {
		s.statsByMinuteCache[minute] = make(map[uint64]database.CounterStats)
	}

	if s.statsByMinuteCache[minute][bannerID].Count == 0 {
		newCounterStats := database.CounterStats{
			BannerID:      bannerID,
			TimestampFrom: uint64(timestampFrom),
			TimestampTo:   uint64(timestampTo),
			Count:         uint64(1),
		}

		s.statsByMinuteCache[minute][bannerID] = newCounterStats
	} else {
		counterStats := s.statsByMinuteCache[minute][bannerID]
		counterStats.Count += uint64(1)
		s.statsByMinuteCache[minute][bannerID] = counterStats
	}
	mutex.Unlock()
}

// Response TODO move it
type Response struct {
	BannerID      uint64 `json:"bannerId"`
	Count         uint64 `json:"count"`
	TimestampFrom uint64 `json:"timestampFrom"`
	TimestampTo   uint64 `json:"timestampTo"`
}

func (s *BannerService) GetBannerCounterStats(bannerID uint64, tsFrom, tsTo uint64) (*Response, error) {
	counterStats, err := s.DBRepo.GetBannerCounterStats(bannerID, tsFrom, tsTo)

	if err != nil {
		return nil, err
	}

	var count uint64

	for _, stat := range counterStats {
		count += stat.Count
	}

	response := &Response{
		BannerID:      bannerID,
		Count:         count,
		TimestampFrom: tsFrom,
		TimestampTo:   tsTo,
	}

	return response, nil
}
