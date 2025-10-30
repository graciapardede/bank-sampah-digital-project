package db

import (
    "fmt"
    "log"
    "os"

    "bank-sampah-digital/backend/models"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDatabase connects to Postgres using environment variables.
// It will run AutoMigrate for all models.
func ConnectDatabase() (*gorm.DB, error) {
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        host := os.Getenv("PGHOST")
        if host == "" {
            host = "localhost"
        }
        port := os.Getenv("PGPORT")
        if port == "" {
            port = "5432"
        }
        user := os.Getenv("PGUSER")
        if user == "" {
            user = "postgres"
        }
        password := os.Getenv("PGPASSWORD")
        dbname := os.Getenv("PGDATABASE")
        if dbname == "" {
            dbname = "bank_sampah"
        }
        sslmode := os.Getenv("PGSSLMODE")
        if sslmode == "" {
            sslmode = "disable"
        }
        dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, password, dbname, port, sslmode)
    }

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    // AutoMigrate all models
    if err := db.AutoMigrate(
        &models.Role{},
        &models.Location{},
        &models.User{},
        &models.WasteType{},
        &models.RewardItem{},
        &models.Deposit{},
        &models.DepositItem{},
        &models.Redemption{},
        &models.RedemptionItem{},
        &models.PointsLedger{},
        &models.AuditLog{},
    ); err != nil {
        log.Printf("AutoMigrate error: %v", err)
        return nil, err
    }

    DB = db
    return db, nil
}
