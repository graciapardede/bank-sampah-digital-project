package main

import (
    "fmt"
    "log"

    "bank-sampah-digital/backend/db"
    "bank-sampah-digital/backend/models"

    "golang.org/x/crypto/bcrypt"
)

func main() {
    d, err := db.ConnectDatabase()
    if err != nil {
        log.Fatalf("connect db: %v", err)
    }

    // create roles
    roles := []models.Role{{Name: "warga"}, {Name: "admin"}, {Name: "superadmin"}}
    for _, r := range roles {
        d.FirstOrCreate(&r, models.Role{Name: r.Name})
    }

    // create a default location
    loc := models.Location{Name: "Cabang Utama", Address: "Alamat cabang utama"}
    d.FirstOrCreate(&loc, models.Location{Name: loc.Name})

    // create admin user if not exists
    var adminRole models.Role
    if err := d.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
        log.Fatalf("admin role missing: %v", err)
    }

    password := "admin123"
    hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    admin := models.User{FullName: "Admin Demo", Email: "admin@local.dev", PasswordHash: string(hash), RoleID: adminRole.ID, LocationID: loc.ID}
    d.FirstOrCreate(&admin, models.User{Email: admin.Email})

    fmt.Println("Seed complete")
    fmt.Printf("Admin credentials: email=%s password=%s\n", admin.Email, password)
}
