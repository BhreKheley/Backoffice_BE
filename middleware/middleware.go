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
	"github.com/jmoiron/sqlx"
)

// Ambil jwtSecret dari environment variable
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// Struct untuk JWT claims yang mencakup role_id dan user_id
type CustomClaims struct {
	UserID int `json:"user_id"`
	RoleID int `json:"role_id"`
	jwt.RegisteredClaims
}

// Middleware untuk validasi token JWT
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

		if err != nil {
			log.Println("JWT Parse Error:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			c.Set("user_id", claims.UserID)
			c.Set("role_id", claims.RoleID)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// Fungsi untuk membuat token JWT dengan role_id
func GenerateToken(userID, roleID int) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		RoleID: roleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Middleware untuk memeriksa permission berdasarkan kode permission
func CheckPermission(permissionCode string, db *sqlx.DB) gin.HandlerFunc {
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
		err := db.Get(&permission, "SELECT * FROM permission WHERE code = $1", permissionCode)
		if err != nil {
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
func hasPermission(roleID, permissionID int, db *sqlx.DB) bool {
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM role_permission WHERE role_id = $1 AND permission_id = $2", roleID, permissionID)
	if err != nil {
		log.Println("Error checking permission:", err)
		return false
	}
	return count > 0
}
