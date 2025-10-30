package models

import "gorm.io/gorm"

type Role struct {
    gorm.Model
    Name string `gorm:"size:50;uniqueIndex;not null" json:"name"`
}
