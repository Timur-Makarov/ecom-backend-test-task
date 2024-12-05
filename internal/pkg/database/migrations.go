package database

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func MigrateDB(db *gorm.DB) error {
	m := gormigrate.New(
		db,
		gormigrate.DefaultOptions,
		GetMigrations(db),
	)

	if err := m.Migrate(); err != nil {
		return err
	}

	return nil
}

func GetMigrations(db *gorm.DB) []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: "Initial_Migration_11_30_2024",
			Migrate: func(tx *gorm.DB) error {
				type Banner struct {
					ID   uint64 `json:"id" gorm:"primaryKey"`
					Name string `json:"name"`
				}

				type CounterStats struct {
					BannerID      uint64 `json:"bannerId"`
					TimestampFrom uint64 `json:"timestampFrom" gorm:"primaryKey"`
					TimestampTo   uint64 `json:"timestampTo" gorm:"primaryKey"`
					Count         uint64 `json:"count"`
				}

				err := tx.AutoMigrate(&Banner{})
				if err != nil {
					return err
				}

				return tx.AutoMigrate(&CounterStats{})
			},
			Rollback: func(tx *gorm.DB) error {
				err := db.Migrator().DropTable("banners")
				if err != nil {
					return err
				}

				return tx.Migrator().DropTable("counter_stats")
			},
		},
	}
}
