package middleware

import (
	"absensi-app/models"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type CustomClaims struct {
	UserID  int    `json:"user_id"`
	RoleID  int    `json:"role_id"`
	Context string `json:"context"`
	jwt.RegisteredClaims
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer ", "", 1))
		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			log.Println("JWT Parse Error:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*CustomClaims); ok {
			c.Set("user_id", claims.UserID)
			c.Set("role_id", claims.RoleID)
			c.Set("context", claims.Context)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func GenerateToken(userID, roleID int, context string) (string, error) {
	claims := CustomClaims{
		UserID:  userID,
		RoleID:  roleID,
		Context: context,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}


// Middleware untuk memeriksa permission berdasarkan kode permission
func CheckPermission(permissionCode string, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID, exists := c.Get("role_id")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Role ID not found"})
			c.Abort()
			return
		}

		roleIDInt, ok := roleID.(int)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid role ID"})
			c.Abort()
			return
		}

		var permission models.Permission
		if err := db.Where("code = ?", permissionCode).First(&permission).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission not found"})
			c.Abort()
			return
		}

		if !hasPermission(roleIDInt, permission.ID, db) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// Fungsi untuk cek apakah role memiliki permission tertentu
func hasPermission(roleID, permissionID int, db *gorm.DB) bool {
	var count int64
	err := db.Table("role_permission").Where("role_id = ? AND permission_id = ?", roleID, permissionID).Count(&count).Error
	if err != nil {
		log.Println("Error checking permission:", err)
		return false
	}
	return count > 0
}
