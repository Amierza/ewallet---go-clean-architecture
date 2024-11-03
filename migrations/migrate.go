package migrations

import (
	"github.com/Amierza/e-wallet/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&entity.User{},
		&entity.TopUp{},
		&entity.Payment{},
		&entity.Transfer{},
	); err != nil {
		return err
	}

	return nil
}
