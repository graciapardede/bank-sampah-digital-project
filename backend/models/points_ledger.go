package models

import "gorm.io/gorm"

type PointsLedger struct {
    gorm.Model
    UserID      uint    `json:"user_id"`
    User        *User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Type        string  `gorm:"size:20;not null" json:"type"` // credit/debit
    PointsAmount float64 `gorm:"not null" json:"points_amount"`
    ReferenceID uint    `json:"reference_id"` // e.g., deposit or redemption id
    BalanceAfter float64 `gorm:"not null" json:"balance_after"`
}
