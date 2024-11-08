package handlers

import (
	"net/http"
	"strconv" // Import untuk konversi tipe
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"absensi-app/models"
)

// GetAllPositions retrieves all positions
func GetAllPositions(c *gin.Context, db *sqlx.DB) {
	var positions []models.Position
	err := db.Select(&positions, "SELECT * FROM position")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, positions)
}

// GetPositionByID retrieves a position by ID
func GetPositionByID(c *gin.Context, db *sqlx.DB) {
	id := c.Param("id")
	var position models.Position
	err := db.Get(&position, "SELECT * FROM position WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Position not found"})
		return
	}
	c.JSON(http.StatusOK, position)
}

// CreatePosition creates a new position
func CreatePosition(c *gin.Context, db *sqlx.DB) {
	var position models.Position
	if err := c.ShouldBindJSON(&position); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	_, err := db.NamedExec("INSERT INTO position (position_name, division_id, created_at, updated_at) VALUES (:position_name, :division_id, NOW(), NOW())", &position)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, position)
}

// UpdatePosition updates an existing position
func UpdatePosition(c *gin.Context, db *sqlx.DB) {
	idStr := c.Param("id") // Ambil ID sebagai string
	id, err := strconv.Atoi(idStr) // Konversi ke int
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var position models.Position
	if err := c.ShouldBindJSON(&position); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	position.ID = id // Assign ke field ID
	_, err = db.NamedExec("UPDATE position SET position_name = :position_name, division_id = :division_id, updated_at = NOW() WHERE id = :id", &position)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, position)
}

// DeletePosition deletes a position by ID
func DeletePosition(c *gin.Context, db *sqlx.DB) {
	idStr := c.Param("id") // Ambil ID sebagai string
	id, err := strconv.Atoi(idStr) // Konversi ke int
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	_, err = db.Exec("DELETE FROM position WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
