package models

import "gorm.io/gorm"

type User struct {
    gorm.Model
    FullName      string   `gorm:"size:150;not null" json:"full_name"`
    Email         string   `gorm:"size:150;uniqueIndex;not null" json:"email"`
    PasswordHash  string   `gorm:"size:255;not null" json:"-"`
    BalancePoints float64  `gorm:"default:0" json:"balance_points"`

    RoleID     uint      `json:"role_id"`
    Role       *Role     `gorm:"foreignKey:RoleID" json:"role,omitempty"`

    LocationID uint       `json:"location_id"`
    Location   *Location  `gorm:"foreignKey:LocationID" json:"location,omitempty"`
}

