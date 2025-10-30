package models

import "gorm.io/gorm"

type Reward struct {
    gorm.Model
    Name        string `gorm:"size:100;not null" json:"name"`
    Description string `gorm:"size:255" json:"description"`
    PointsCost  int    `gorm:"not null" json:"points_cost"`
}
