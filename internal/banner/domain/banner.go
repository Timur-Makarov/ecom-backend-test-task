package domain

// The idea is that each layer should have it's own domain models
// so that the layers below would have to adjust themselves.
// For now, there is no usage for them.

type Banner struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type CounterStats struct {
	BannerID      uint64 `json:"bannerId"`
	TimestampFrom uint64 `json:"timestampFrom"`
	TimestampTo   uint64 `json:"timestampTo"`
	Count         uint64 `json:"count"`
}
