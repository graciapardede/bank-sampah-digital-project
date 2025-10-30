package handlers

import (
    "net/http"
    "strconv"

    "bank-sampah-digital/backend/db"
    "bank-sampah-digital/backend/models"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "gorm.io/gorm/clause"
)

// GetPendingDeposits returns deposits with status pending limited to the admin's location
func GetPendingDeposits(c *gin.Context) {
    lidI, exists := c.Get("userLocationID")
    if !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "user location not found in token"})
        return
    }
    locationID := lidI.(uint)

    var deposits []models.Deposit
    if err := db.DB.Preload("Items").Where("status = ? AND location_id = ?", "pending", locationID).Find(&deposits).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, deposits)
}

// VerifyDeposit verifies a deposit and credits points to the user within a transaction
func VerifyDeposit(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    // Start transaction
    err = db.DB.Transaction(func(tx *gorm.DB) error {
        var dep models.Deposit
        if err := tx.Preload("Items").First(&dep, uint(id)).Error; err != nil {
            return err
        }

        if dep.Status != "pending" {
            return nil // nothing to do
        }

        // Calculate total points from items if needed
        total := dep.TotalPoints
        if total == 0 {
            // try compute from items and waste types
            var sum float64
            for _, it := range dep.Items {
                var wt models.WasteType
                if err := tx.First(&wt, it.WasteTypeID).Error; err != nil {
                    return err
                }
                sum += it.Weight * wt.PointsPerUnit
            }
            total = sum
        }

        // Update deposit status
        if err := tx.Model(&dep).Update("status", "verified").Error; err != nil {
            return err
        }

        // Lock user row for update
        var user models.User
        if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, dep.UserID).Error; err != nil {
            return err
        }

        user.BalancePoints += total
        if err := tx.Save(&user).Error; err != nil {
            return err
        }

        // create ledger entry
        ledger := models.PointsLedger{
            UserID:      user.ID,
            Type:        "credit",
            PointsAmount: total,
            ReferenceID: dep.ID,
            BalanceAfter: user.BalancePoints,
        }
        if err := tx.Create(&ledger).Error; err != nil {
            return err
        }

        return nil
    })

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "deposit verified"})
}
