package handlers

import (
	"absensi-app/models"
	"database/sql"
	"fmt"
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
		"api_token":   tokenString,
	})
}

func GetUserByToken(c *gin.Context, db *sqlx.DB) {
	// Mendapatkan token dari header Authorization
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		return
	}

	// Parsing token JWT
	tokenString := authHeader[len("Bearer "):]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validasi algoritma token
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Mengextract klaim dari token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	userID := claims["user_id"].(float64) // Mengambil user ID dari klaim

	// Query user dari database berdasarkan user ID
	var user models.User
	err = db.QueryRowx(`SELECT id, username, email, role_id, is_active FROM "user" WHERE id = $1`, int(userID)).StructScan(&user)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Jika user ditemukan, return data user
	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role_id":  user.RoleID,
		"is_active": user.IsActive,
	})
}
