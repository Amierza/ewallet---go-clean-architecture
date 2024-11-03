package entity

import "github.com/google/uuid"

type Payment struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"payment_id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	User          User      `gorm:"foreignKey:UserID"`
	Amount        int64     `json:"amount"`
	Remarks       string    `gorm:"type:text;null" json:"remarks"`
	BalanceBefore int64     `json:"balance_before"`
	BalanceAfter  int64     `json:"balance_after"`
	Timestamp
}
