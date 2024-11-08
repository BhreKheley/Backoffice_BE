package handlers

import (
	"net/http"
	"strconv" // Import untuk konversi tipe
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"absensi-app/models"
)

// GetAllDivisions retrieves all divisions
func GetAllDivisions(c *gin.Context, db *sqlx.DB) {
	var divisions []models.Division
	err := db.Select(&divisions, "SELECT * FROM division")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, divisions)
}

// GetDivisionByID retrieves a division by ID
func GetDivisionByID(c *gin.Context, db *sqlx.DB) {
	id := c.Param("id")
	var division models.Division
	err := db.Get(&division, "SELECT * FROM division WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Division not found"})
		return
	}
	c.JSON(http.StatusOK, division)
}

// CreateDivision creates a new division
func CreateDivision(c *gin.Context, db *sqlx.DB) {
	var division models.Division
	if err := c.ShouldBindJSON(&division); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	_, err := db.NamedExec("INSERT INTO division (division_name, created_at, updated_at) VALUES (:division_name, NOW(), NOW())", &division)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, division)
}

// UpdateDivision updates an existing division
func UpdateDivision(c *gin.Context, db *sqlx.DB) {
	idStr := c.Param("id") // Ambil ID sebagai string
	id, err := strconv.Atoi(idStr) // Konversi ke int
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var division models.Division
	if err := c.ShouldBindJSON(&division); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	division.ID = id // Assign ke field ID
	_, err = db.NamedExec("UPDATE division SET division_name = :division_name, updated_at = NOW() WHERE id = :id", &division)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, division)
}

// DeleteDivision deletes a division by ID
func DeleteDivision(c *gin.Context, db *sqlx.DB) {
	idStr := c.Param("id") // Ambil ID sebagai string
	id, err := strconv.Atoi(idStr) // Konversi ke int
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	_, err = db.Exec("DELETE FROM division WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
