package handlers

import (
	"absensi-app/models"
	"database/sql"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

var jwtSecret = []byte("secretkey") // Lu bisa ganti ini dengan key yang lebih aman

func Login(c *gin.Context, db *sqlx.DB) {
	var loginData models.User
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	// Query user berdasarkan email
	err := db.QueryRow(`SELECT id, username, email, password, role_id, is_active FROM "user" WHERE email = $1`, loginData.Email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.RoleID, &user.IsActive)

	if err == sql.ErrNoRows {
		// Jika user tidak ditemukan
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	} else if err != nil {
		// Jika ada error lain
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Compare hashed password dari database dengan password yang diinput
	// err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	// if err != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
	// 	return
	// }

	// Cek apakah user aktif
	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User account is inactive"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role_id": user.RoleID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token valid 24 jam
	})

	// Sign token menggunakan secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Login sukses, return token dan info user (tanpa password)
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful!",
		"token":   tokenString,
		"user": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"role_id":   user.RoleID,
			"is_active": user.IsActive,
		},
	})
}
