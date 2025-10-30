package middleware

import (
    "net/http"
    "os"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
)

// AuthMiddleware validates JWT and sets userID and userLocationID in context
func AuthMiddleware() gin.HandlerFunc {
    secret := []byte(getEnv("JWT_SECRET", "secret"))
    return func(c *gin.Context) {
        ah := c.GetHeader("Authorization")
        if ah == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
            return
        }
        parts := strings.SplitN(ah, " ", 2)
        if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization header"})
            return
        }
        tokenStr := parts[1]
        token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
            if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, jwt.ErrTokenMalformed
            }
            return secret, nil
        })
        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            return
        }
        if claims, ok := token.Claims.(jwt.MapClaims); ok {
            // Extract userID and locationID if present
            if uid, ok := claims["user_id"]; ok {
                c.Set("userID", uint(toUint64(uid)))
            }
            if lid, ok := claims["location_id"]; ok {
                c.Set("userLocationID", uint(toUint64(lid)))
            }
            if role, ok := claims["role"]; ok {
                if s, ok := role.(string); ok {
                    c.Set("role", s)
                }
            }
        }
        c.Next()
    }
}

// RoleMiddleware ensures the user has the required role
func RoleMiddleware(roleName string) gin.HandlerFunc {
    return func(c *gin.Context) {
        r, exists := c.Get("role")
        if !exists {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "role not found"})
            return
        }
        if roleStr, ok := r.(string); ok {
            if roleStr != roleName {
                c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient role"})
                return
            }
        }
        c.Next()
    }
}

func getEnv(key, def string) string {
    v := os.Getenv(key)
    if v == "" {
        return def
    }
    return v
}

// helper to cast jwt claim numeric to uint64
func toUint64(v interface{}) uint64 {
    switch t := v.(type) {
    case float64:
        return uint64(t)
    case float32:
        return uint64(t)
    case int:
        return uint64(t)
    case int64:
        return uint64(t)
    case uint64:
        return t
    case string:
        // try parse
        var out uint64
        for i := 0; i < len(t); i++ {
            out = out*10 + uint64(t[i]-'0')
        }
        return out
    default:
        return 0
    }
}
