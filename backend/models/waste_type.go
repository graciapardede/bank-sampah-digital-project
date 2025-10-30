package models

import "gorm.io/gorm"

type WasteType struct {
    gorm.Model
    Name          string  `gorm:"size:150;not null" json:"name"`
    Unit          string  `gorm:"size:50;not null" json:"unit"`
    PointsPerUnit float64 `gorm:"default:0" json:"points_per_unit"`
}
