package controllers

import "gorm.io/gorm"

var DB *gorm.DB

// SetDB sets the package-level DB pointer used by controllers
func SetDB(db *gorm.DB) {
    DB = db
}
