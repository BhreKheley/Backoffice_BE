package handlers

import (
	"absensi-app/helpers"
	"absensi-app/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func GetUserByID(c *gin.Context, db *sqlx.DB) {
	id := c.Param("id")

	var user models.User
	err := db.QueryRow("SELECT id, username, email, role_id, is_active FROM \"user\" WHERE id = $1", id).Scan(&user.ID, &user.Username, &user.Email, &user.RoleID, &user.IsActive)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func GetAllUsers(c *gin.Context, db *sqlx.DB) {
	var user []models.User
	err := db.Select(&user, `SELECT id, username, email, role_id, is_active FROM "user"`)

	if err != nil {
		fmt.Println("Error: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateUser handles the creation of a new user with a hashed password
func CreateUser(c *gin.Context, db *sqlx.DB) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the user's password
	hashedPassword, err := helpers.HashPassword(newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Insert the new user into the database
	_, err = db.Exec(`INSERT INTO "user" (username, email, password, role_id, is_active) 
	                  VALUES ($1, $2, $3, $4, $5)`,
		newUser.Username, newUser.Email, hashedPassword, newUser.RoleID, newUser.IsActive)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully!"})
}
