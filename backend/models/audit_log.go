package models

import (
    "time"
    "gorm.io/gorm"
)

type AuditLog struct {
    gorm.Model
    Action    string    `gorm:"size:150;not null" json:"action"`
    UserID    *uint     `json:"user_id"`
    Details   string    `gorm:"size:1000" json:"details"`
    CreatedAt time.Time `json:"created_at"`
}
