package entity

import "github.com/google/uuid"

type TopUp struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"top_up_id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	User          User      `gorm:"foreignKey:UserID"`
	Amount        int64     `json:"amount"`
	BalanceBefore int64     `json:"balance_before"`
	BalanceAfter  int64     `json:"balance_after"`
	Timestamp
}
