package handlers

import (
	"absensi-app/helpers"
	"absensi-app/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// GetUserByID handles getting a user by ID
func GetUserByID(c *gin.Context, db *sqlx.DB) {
	id := c.Param("id")

	var user models.User
	err := db.QueryRow("SELECT id, username, email, role_id, is_active FROM \"user\" WHERE id = $1", id).
		Scan(&user.ID, &user.Username, &user.Email, &user.RoleID, &user.IsActive)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"role_id":   user.RoleID,
			"is_active": user.IsActive,
		},
	})
}

// GetAllUsers handles getting all users
func GetAllUsers(c *gin.Context, db *sqlx.DB) {
	var users []models.User
	err := db.Select(&users, `SELECT id, username, email, role_id, is_active FROM "user"`)

	if err != nil {
		fmt.Println("Error: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to retrieve users",
		})
		return
	}

	// If no users found
	if len(users) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "No users found",
			"data":    []interface{}{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   formatUsers(users), // Helper function to format the result
	})
}

// CreateUser handles the creation of a new user
func CreateUser(c *gin.Context, db *sqlx.DB) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"error":   err.Error(),
		})
		return
	}

	// Hash the user's password
	hashedPassword, err := helpers.HashPassword(newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to hash password",
		})
		return
	}

	// Insert the new user into the database
	_, err = db.Exec(`INSERT INTO "user" (username, email, password, role_id, is_active) 
	                  VALUES ($1, $2, $3, $4, $5)`,
		newUser.Username, newUser.Email, hashedPassword, newUser.RoleID, newUser.IsActive)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to create user",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User created successfully!",
		"data": gin.H{
			"username": newUser.Username,
			"email":    newUser.Email,
			"role_id":  newUser.RoleID,
		},
	})
}

// UpdateUser handles updating an existing user by ID
func UpdateUser(c *gin.Context, db *sqlx.DB) {
	id := c.Param("id")
	var updatedUser models.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"error":   err.Error(),
		})
		return
	}

	// Update the user in the database
	_, err := db.Exec(`UPDATE "user" SET username = $1, email = $2, role_id = $3, is_active = $4 WHERE id = $5`,
		updatedUser.Username, updatedUser.Email, updatedUser.RoleID, updatedUser.IsActive, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to update user",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User updated successfully!",
	})
}

// DeleteUser handles deleting a user by ID
func DeleteUser(c *gin.Context, db *sqlx.DB) {
	id := c.Param("id")

	// Delete the user from the database
	_, err := db.Exec(`DELETE FROM "user" WHERE id = $1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to delete user",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User deleted successfully!",
	})
}

// Helper function to format users
func formatUsers(users []models.User) []gin.H {
	var formattedUsers []gin.H
	for _, user := range users {
		formattedUsers = append(formattedUsers, gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"role_id":   user.RoleID,
			"is_active": user.IsActive,
		})
	}
	return formattedUsers
}
