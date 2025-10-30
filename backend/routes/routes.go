package routes

import (
    "bank-sampah-digital/backend/handlers"
    "bank-sampah-digital/backend/middleware"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
    // Public auth routes
    r.POST("/auth/register", handlers.Register)
    r.POST("/auth/login", handlers.Login)

    // Admin protected routes
    admin := r.Group("/admin")
    admin.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
    {
        admin.GET("/deposits/pending", handlers.GetPendingDeposits)
        admin.POST("/deposits/verify/:id", handlers.VerifyDeposit)
        admin.POST("/redemptions/confirm/:id", handlers.ConfirmRedemption)
    }
}
