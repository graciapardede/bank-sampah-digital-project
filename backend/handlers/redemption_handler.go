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

// ConfirmRedemption confirms a redemption: marks confirmed, debits user and reduces reward stock
func ConfirmRedemption(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    err = db.DB.Transaction(func(tx *gorm.DB) error {
        var red models.Redemption
        if err := tx.Preload("Items").First(&red, uint(id)).Error; err != nil {
            return err
        }

        if red.Status != "pending" {
            return nil
        }

        // compute total points if needed
        total := red.TotalPoints
        if total == 0 {
            var sum float64
            for _, it := range red.Items {
                var ri models.RewardItem
                if err := tx.First(&ri, it.RewardItemID).Error; err != nil {
                    return err
                }
                sum += float64(it.Quantity) * float64(ri.PointsCost)
            }
            total = sum
        }

        // lock user
        var user models.User
        if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, red.UserID).Error; err != nil {
            return err
        }

        if user.BalancePoints < total {
            return gorm.ErrInvalidData
        }

        // reduce user balance
        user.BalancePoints -= total
        if err := tx.Save(&user).Error; err != nil {
            return err
        }

        // create ledger
        ledger := models.PointsLedger{
            UserID:      user.ID,
            Type:        "debit",
            PointsAmount: total,
            ReferenceID: red.ID,
            BalanceAfter: user.BalancePoints,
        }
        if err := tx.Create(&ledger).Error; err != nil {
            return err
        }

        // reduce stock for each item but ensure location matches
        for _, it := range red.Items {
            var ri models.RewardItem
            if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND location_id = ?", it.RewardItemID, red.LocationID).First(&ri).Error; err != nil {
                return err
            }
            if ri.Stock < it.Quantity {
                return gorm.ErrInvalidData
            }
            ri.Stock -= it.Quantity
            if err := tx.Save(&ri).Error; err != nil {
                return err
            }
        }

        // update redemption status
        if err := tx.Model(&red).Update("status", "confirmed").Error; err != nil {
            return err
        }

        return nil
    })

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "redemption confirmed"})
}
