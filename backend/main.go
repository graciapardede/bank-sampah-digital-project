package main

import (
    "log"

    "bank-sampah-digital/backend/db"
    "bank-sampah-digital/backend/handlers"
    "bank-sampah-digital/backend/middleware"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

func main() {
    // connect to DB
    _, err := db.ConnectDatabase()
    if err != nil {
        log.Fatalf("failed to connect db: %v", err)
    }

    r := gin.Default()

    // Enable CORS for local frontend dev
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "Bank Sampah Digital API"})
    })

    adminGroup := r.Group("/admin")
    adminGroup.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
    {
        adminGroup.GET("/deposits/pending", handlers.GetPendingDeposits)
        adminGroup.POST("/deposits/verify/:id", handlers.VerifyDeposit)
        adminGroup.POST("/redemptions/confirm/:id", handlers.ConfirmRedemption)
    }

    log.Println("Server running on :8080")
    if err := r.Run(":8080"); err != nil {
        log.Fatalf("server exited: %v", err)
    }
}
