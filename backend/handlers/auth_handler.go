package handlers

import (
    "net/http"
    "os"
    "time"

    "bank-sampah-digital/backend/db"
    "bank-sampah-digital/backend/models"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
    "golang.org/x/crypto/bcrypt"
)

type registerReq struct {
    FullName   string `json:"full_name" binding:"required"`
    Email      string `json:"email" binding:"required,email"`
    Password   string `json:"password" binding:"required,min=6"`
    RoleID     uint   `json:"role_id"`
    LocationID uint   `json:"location_id"`
}

type loginReq struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

func getJWTSecret() []byte {
    s := os.Getenv("JWT_SECRET")
    if s == "" {
        s = "secret"
    }
    return []byte(s)
}

func createToken(user *models.User, roleName string) (string, error) {
    claims := jwt.MapClaims{
        "user_id":     user.ID,
        "role":        roleName,
        "location_id": user.LocationID,
        "exp":         time.Now().Add(7 * 24 * time.Hour).Unix(),
    }
    t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return t.SignedString(getJWTSecret())
}

// Register creates a new user
func Register(c *gin.Context) {
    var req registerReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
        return
    }

    user := models.User{
        FullName:     req.FullName,
        Email:        req.Email,
        PasswordHash: string(hashed),
        RoleID:       req.RoleID,
        LocationID:   req.LocationID,
    }

    if err := db.DB.Create(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // load role name if possible
    var role models.Role
    roleName := "warga"
    if user.RoleID != 0 {
        if err := db.DB.First(&role, user.RoleID).Error; err == nil {
            roleName = role.Name
        }
    }

    token, err := createToken(&user, roleName)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"token": token, "user": gin.H{"id": user.ID, "email": user.Email, "full_name": user.FullName, "role": roleName, "location_id": user.LocationID}})
}

// Login authenticates and returns a JWT
func Login(c *gin.Context) {
    var req loginReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user models.User
    if err := db.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }

    var role models.Role
    roleName := "warga"
    if user.RoleID != 0 {
        if err := db.DB.First(&role, user.RoleID).Error; err == nil {
            roleName = role.Name
        }
    }

    token, err := createToken(&user, roleName)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token, "user": gin.H{"id": user.ID, "email": user.Email, "full_name": user.FullName, "role": roleName, "location_id": user.LocationID}})
}
