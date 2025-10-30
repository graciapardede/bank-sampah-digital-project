package models

import "gorm.io/gorm"

type Location struct {
    gorm.Model
    Name    string `gorm:"size:150;not null" json:"name"`
    Address string `gorm:"size:300" json:"address"`
}
