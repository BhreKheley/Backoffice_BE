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

	c.JSON(http.StatusOK, LoginResponse{
		Message:  "Login berhasil!",
		APIToken: tokenString,
	})
}

type GetUserByTokenResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	RoleID   int    `json:"role_id"`
	IsActive bool   `json:"is_active"`
}

func GetUserByToken(c *gin.Context, db *gorm.DB) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User ID tidak ditemukan di konteks"})
		return
	}

	userIDInt, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Tipe user ID tidak valid"})
		return
	}

	var user models.User
	if err := db.First(&user, userIDInt).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Pengguna tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Gagal mengakses database"})
		return
	}

	c.JSON(http.StatusOK, GetUserByTokenResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		RoleID:   user.RoleID,
		IsActive: user.IsActive,
	})
}
