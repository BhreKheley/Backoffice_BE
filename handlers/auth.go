package handlers

import (
	"absensi-app/middleware"
	"absensi-app/models"
	"net/http"
	// "time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// Login handler untuk otentikasi pengguna
func Login(c *gin.Context, db *sqlx.DB) {
	var loginData models.User
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	// Query user berdasarkan email
	err := db.QueryRowx(`SELECT id, username, email, password, role_id, is_active FROM "user" WHERE email = $1`, loginData.Email).
		StructScan(&user)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Cek apakah user aktif
	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User account is inactive"})
		return
	}

	// Cek password
	if user.Password != loginData.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT token dengan role_id
	tokenString, err := middleware.GenerateToken(user.ID, user.RoleID) // Gunakan GenerateToken dari middleware
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Login sukses, return token dan info user (tanpa password)
	c.JSON(http.StatusOK, gin.H{
		"message":   "Login successful!",
		"api_token": tokenString,
	})
}

// GetUserByToken handler untuk mendapatkan data user dari token
func GetUserByToken(c *gin.Context, db *sqlx.DB) {
	// Ambil user_id dari context yang sudah diset di middleware AuthMiddleware
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	// Konversi userID ke integer
	userIDInt, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID type"})
		return
	}

	// Query user dari database berdasarkan user ID
	var user models.User
	err := db.QueryRowx(`SELECT id, username, email, role_id, is_active FROM "user" WHERE id = $1`, userIDInt).StructScan(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Jika user ditemukan, return data user
	c.JSON(http.StatusOK, gin.H{
		"id":        user.ID,
		"username":  user.Username,
		"email":     user.Email,
		"role_id":   user.RoleID,
		"is_active": user.IsActive,
	})
}
