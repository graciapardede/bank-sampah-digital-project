package middleware

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

// RequireRole returns a middleware that checks the role in context
func RequireRole(allowed string) gin.HandlerFunc {
    return func(c *gin.Context) {
        r, exists := c.Get("role")
        if !exists {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "role not found"})
            return
        }
        roleStr, _ := r.(string)
        if roleStr != allowed {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient role"})
            return
        }
        c.Next()
    }
}
