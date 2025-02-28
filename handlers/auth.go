package handlers

import (
	"absensi-app/middleware"
	"absensi-app/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginResponse struct {
	Message  string `json:"message"`
	APIToken string `json:"api_token"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func Login(c *gin.Context, db *gorm.DB, context string) {
	var loginData struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Input tidak valid"})
		return
	}

	var user models.User
	if err := db.Where("email = ?", loginData.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Email atau password tidak valid"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Gagal mengakses database"})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Akun pengguna tidak aktif"})
		return
	}

	if user.Password != loginData.Password {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Email atau password tidak valid"})
		return
	}

	tokenString, err := middleware.GenerateToken(user.ID, user.RoleID, context)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Gagal membuat token"})
		return
	}

	// Simpan token ke database
	if err := db.Model(&user).Update("current_token", tokenString).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Gagal menyimpan token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Message:  "Login berhasil!",
		APIToken: tokenString,
	})
}

type GetUserByTokenResponse struct {
	UserID       int    `json:"user_id"`
	EmployeeID   int    `json:"employee_id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	RoleID       int    `json:"role_id"`
	IsActive     bool   `json:"is_active"`
	FullName     string `json:"full_name"`
	Avatar       string `json:"avatar"`
	Phone        string `json:"phone"`
	PositionID   int    `json:"position_id"`
	DivisionID   int    `json:"division_id"`
	PositionName string `json:"position_name"`
	DivisionName string `json:"division_name"`
}

// GetUserByToken - Mengecek token dan mengambil data pengguna
func GetUserByToken(c *gin.Context, db *gorm.DB) {
	// Cek apakah token valid
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "User ID tidak ditemukan, silakan login kembali"})
		return
	}

	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Token tidak ditemukan, silakan login kembali"})
		return
	}

	tokenString = middleware.ExtractBearerToken(tokenString)

	// Ambil user dari database berdasarkan token
	var user models.User
	if err := db.Where("id = ? AND current_token = ?", userID, tokenString).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Token tidak valid atau sudah expired, silakan login kembali"})
		return
	}

	// Ambil data user + employee + posisi + divisi
	var response GetUserByTokenResponse
	if err := db.Table(`"user"`).
		Select(`"user".id AS user_id, "user".username, "user".email, "user".role_id, "user".is_active, 
				e.id AS employee_id, e.full_name, e.avatar, e.phone, e.position_id, e.division_id, 
				p.position_name, d.division_name`).
		Joins(`LEFT JOIN employee AS e ON "user".id = e.user_id`).
		Joins("LEFT JOIN position AS p ON e.position_id = p.id").
		Joins("LEFT JOIN division AS d ON e.division_id = d.id").
		Where(`"user".id = ?`, userID).
		Scan(&response).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengambil data pengguna"})
		return
	}

	// Token valid, kembalikan informasi user
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Token masih valid",
		"data":    response,
	})
}
