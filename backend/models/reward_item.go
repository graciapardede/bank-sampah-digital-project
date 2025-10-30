package models

import "gorm.io/gorm"

type RewardItem struct {
    gorm.Model
    Name       string   `gorm:"size:150;not null" json:"name"`
    PointsCost int      `gorm:"not null" json:"points_cost"`
    Stock      int      `gorm:"default:0" json:"stock"`
    LocationID uint     `json:"location_id"`
    Location   *Location `gorm:"foreignKey:LocationID" json:"location,omitempty"`
}
