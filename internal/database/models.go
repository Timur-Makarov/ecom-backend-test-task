package database

type Banner struct {
	ID   uint64 `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

// CounterStats is a table with per-minute (per any timestamp range) user interactions with banners.
// Unix timestamp is int64, but uint64 is used to show that we only considering events after the epoch.
type CounterStats struct {
	BannerID      uint64 `json:"bannerId"`
	TimestampFrom uint64 `json:"timestampFrom" gorm:"primaryKey"`
	TimestampTo   uint64 `json:"timestampTo" gorm:"primaryKey"`
	Count         uint64 `json:"count"`
}
