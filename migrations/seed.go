package migrations

import (
	"github.com/Amierza/e-wallet/migrations/seeds"
	"gorm.io/gorm"
)

func Seeder(db *gorm.DB) error {
	if err := seeds.ListUserSeeder(db); err != nil {
		return err
	}
	return nil
}
