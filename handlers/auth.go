package handlers

import (
	"absensi-app/middleware"
	"absensi-app/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LoginResponse represents the response structure for login
type LoginResponse struct {
	Message   string `json:"message"`
	APIToken  string `json:"api_token"`
}

// ErrorResponse represents the structure for error responses
type ErrorResponse struct {
	Error string `json:"error"`
}

// Login handler untuk otentikasi pengguna
func Login(c *gin.Context, db *gorm.DB) {
	var loginData struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Bind JSON data to loginData struct
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Input tidak valid"})
		return
	}

	var user models.User
	// Query user berdasarkan email menggunakan GORM
	if err := db.Where("email = ?", loginData.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Email atau password tidak valid"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Gagal mengakses database"})
		return
	}

	// Cek apakah user aktif
	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Akun pengguna tidak aktif"})
		return
	}

	// Cek password (pastikan hashing password digunakan)
	if user.Password != loginData.Password {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Email atau password tidak valid"})
		return
	}

	// Generate JWT token dengan role_id
	tokenString, err := middleware.GenerateToken(user.ID, user.RoleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Gagal membuat token"})
		return
	}

	// Login sukses, return token dan info user (tanpa password)
	c.JSON(http.StatusOK, LoginResponse{
		Message:  "Login berhasil!",
		APIToken: tokenString,
	})
}

// GetUserByTokenResponse represents the response structure for user data retrieved from token
type GetUserByTokenResponse struct {
	ID        int   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	RoleID    int   `json:"role_id"`
	IsActive  bool   `json:"is_active"`
}

// GetUserByToken handler untuk mendapatkan data user dari token
func GetUserByToken(c *gin.Context, db *gorm.DB) {
	// Ambil user_id dari context yang sudah diset di middleware AuthMiddleware
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User ID tidak ditemukan di konteks"})
		return
	}

	// Konversi userID ke uint
	userIDInt, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Tipe user ID tidak valid"})
		return
	}

	// Query user dari database berdasarkan user ID menggunakan GORM
	var user models.User
	if err := db.First(&user, userIDInt).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Pengguna tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Gagal mengakses database"})
		return
	}

	// Jika user ditemukan, return data user
	c.JSON(http.StatusOK, GetUserByTokenResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		RoleID:   user.RoleID,
		IsActive: user.IsActive,
	})
}
