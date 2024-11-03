package entity

import "github.com/google/uuid"

type Transfer struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"transfer_id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	TargetUserID  uuid.UUID `gorm:"type:uuid;not null" json:"target_user_id"`
	User          User      `gorm:"foreignKey:UserID"`
	TargetUser    User      `gorm:"foreignKey:TargetUserID;references:ID"`
	Amount        int64     `json:"amount"`
	Remarks       string    `gorm:"type:text;null" json:"remarks"`
	BalanceBefore int64     `json:"balance_before"`
	BalanceAfter  int64     `json:"balance_after"`
	Timestamp
}
