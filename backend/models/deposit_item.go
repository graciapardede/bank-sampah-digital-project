package models

import "gorm.io/gorm"

type DepositItem struct {
    gorm.Model
    DepositID   uint      `json:"deposit_id"`
    Deposit     *Deposit  `gorm:"foreignKey:DepositID" json:"deposit,omitempty"`
    WasteTypeID uint      `json:"waste_type_id"`
    WasteType   *WasteType `gorm:"foreignKey:WasteTypeID" json:"waste_type,omitempty"`
    Weight      float64   `json:"weight"`
}
