package pgg

type Banner struct {
	ID   int32  `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

type CounterStatistic struct {
	BannerID      int32 `json:"bannerId"`
	TimestampFrom int64 `json:"timestampFrom" gorm:"primaryKey"`
	TimestampTo   int64 `json:"timestampTo" gorm:"primaryKey"`
	Count         int64 `json:"count"`
}
