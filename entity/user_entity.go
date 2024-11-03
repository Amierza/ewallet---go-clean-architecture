package entity

import (
	"github.com/Amierza/e-wallet/helpers"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey" json:"user_id"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	PhoneNumber string     `json:"phone_number"`
	Address     string     `json:"address"`
	Pin         string     `json:"pin"`
	Balance     int64      `json:"balance"`
	TopUps      []TopUp    `gorm:"foreignKey:UserID"`
	Payments    []Payment  `gorm:"foreignKey:UserID"`
	Transfers   []Transfer `gorm:"foreignKey:UserID"`
	Timestamp
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	u.ID = uuid.New()

	var err error
	u.Pin, err = helpers.HashPin(u.Pin)
	if err != nil {
		return err
	}
	return nil
}
