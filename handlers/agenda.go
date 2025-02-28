package handlers

import (
	"absensi-app/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler untuk membuat agenda (hanya admin/HR)
func CreateAgenda(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var agenda models.Agenda
		if err := c.ShouldBindJSON(&agenda); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
			return
		}

		if err := db.Create(&agenda).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat agenda"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Agenda berhasil dibuat", "data": agenda})
	}
}

// Handler untuk mendapatkan semua agenda
func GetAllAgendas(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var agendas []models.Agenda
		if err := db.Find(&agendas).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data agenda"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": agendas})
	}
}

// Handler untuk mendapatkan detail agenda berdasarkan ID
func GetAgendaDetail(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID agenda tidak valid"})
			return
		}

		var agenda models.Agenda
		if err := db.Preload("Creator").First(&agenda, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Agenda tidak ditemukan"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data agenda"})
			}
			return
		}

		// Ambil daftar peserta
		var participants []models.AgendaParticipant
		if err := db.Preload("User").Where("agenda_id = ?", id).Find(&participants).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil peserta"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"agenda": agenda, "participants": participants})
	}
}

// Handler untuk memperbarui agenda
func UpdateAgenda(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID agenda tidak valid"})
			return
		}

		var agenda models.Agenda
		if err := db.First(&agenda, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Agenda tidak ditemukan"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data agenda"})
			}
			return
		}

		if err := c.ShouldBindJSON(&agenda); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
			return
		}

		if err := db.Save(&agenda).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui agenda"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Agenda berhasil diperbarui", "data": agenda})
	}
}

// Handler untuk menghapus agenda
func DeleteAgenda(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID agenda tidak valid"})
			return
		}

		if err := db.Delete(&models.Agenda{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus agenda"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Agenda berhasil dihapus"})
	}
}

// Handler untuk user menambahkan partisipasi ke agenda
func JoinAgenda(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.Atoi(c.GetString("user_id"))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak valid"})
			return
		}

		agendaID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID agenda tidak valid"})
			return
		}

		// Cek apakah user sudah terdaftar
		var existing models.AgendaParticipant
		if err := db.Where("agenda_id = ? AND user_id = ?", agendaID, userID).First(&existing).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Anda sudah terdaftar dalam agenda ini"})
			return
		}

		participant := models.AgendaParticipant{
			AgendaID: uint(agendaID),
			UserID:   uint(userID),
			AddedBy:  uint(userID),
		}

		if err := db.Create(&participant).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendaftar ke agenda"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Partisipasi berhasil ditambahkan, menunggu persetujuan"})
	}
}

// Handler untuk HR/Admin menambahkan user ke agenda
func AddParticipant(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminID, err := strconv.Atoi(c.GetString("user_id"))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Admin tidak valid"})
			return
		}

		agendaID, err := strconv.Atoi(c.Param("agenda_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID agenda tidak valid"})
			return
		}

		userID, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID user tidak valid"})
			return
		}

		// Cek apakah user sudah ada di agenda
		var existing models.AgendaParticipant
		if err := db.Where("agenda_id = ? AND user_id = ?", agendaID, userID).First(&existing).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User sudah terdaftar dalam agenda ini"})
			return
		}

		participant := models.AgendaParticipant{
			AgendaID: uint(agendaID),
			UserID:   uint(userID),
			AddedBy:  uint(adminID),
		}

		if err := db.Create(&participant).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan peserta"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Peserta berhasil ditambahkan"})
	}
}

// Handler untuk HR/Admin menghapus peserta dari agenda
func RemoveParticipant(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		agendaID, err := strconv.Atoi(c.Param("agenda_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID agenda tidak valid"})
			return
		}

		userID, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID user tidak valid"})
			return
		}

		if err := db.Where("agenda_id = ? AND user_id = ?", agendaID, userID).Delete(&models.AgendaParticipant{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus peserta"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Peserta berhasil dihapus"})
	}
}
