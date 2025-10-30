package models

import "gorm.io/gorm"

type Deposit struct {
    gorm.Model
    UserID      uint         `json:"user_id"`
    User        *User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
    AdminID     *uint        `json:"admin_id"`
    Admin       *User        `gorm:"foreignKey:AdminID" json:"admin,omitempty"`
    LocationID  uint         `json:"location_id"`
    Location    *Location    `gorm:"foreignKey:LocationID" json:"location,omitempty"`
    TotalPoints float64      `gorm:"default:0" json:"total_points"`
    Status      string       `gorm:"size:50;default:pending" json:"status"`
    Items       []DepositItem `gorm:"foreignKey:DepositID" json:"items,omitempty"`
}
