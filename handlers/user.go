package handlers

import (
	"absensi-app/helpers"
	"absensi-app/models"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserResponse struct untuk response detail user
type UserResponse struct {
	Status string     `json:"status"`
	Data   UserDetail `json:"data"`
}

// UserDetail struct untuk detail user
type UserDetail struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	RoleID   int    `json:"role_id"`
	RoleName string `json:"role_name"`
	IsActive bool   `json:"is_active"`
}

// UserListResponse struct untuk response list user
type UserListResponse struct {
	Status string       `json:"status"`
	Data   []UserDetail `json:"data"`
}

// GetUserByID validation
func GetUserByID(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var user struct {
		models.User
		RoleName string `gorm:"column:role_name"`
	}

	err := db.Table(`"user" AS u`).
		Select("u.id, u.username, u.email, u.role_id, u.is_active, COALESCE(r.role_name, 'Unknown') AS role_name").
		Joins("LEFT JOIN role AS r ON u.role_id = r.id").
		Where("u.id = ?", id).
		First(&user).Error // Changed from Scan to First for proper record existence check

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "User not found"})
		} else {
			fmt.Println("Detailed Error: ", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve user"})
		}
		return
	}

	// If the user ID is zero (default value), it indicates no record was found
	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		Status: "success",
		Data: UserDetail{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			RoleID:   user.RoleID,
			RoleName: user.RoleName,
			IsActive: user.IsActive,
		},
	})
}

// GetAllUsers handles fetching all users
func GetAllUsers(c *gin.Context, db *gorm.DB) {
	var user []struct {
		models.User
		RoleName string `gorm:"column:role_name"`
	}

	err := db.Table(`"user" AS u`).
		Select("u.id, u.username, u.email, u.role_id, u.is_active, COALESCE(r.role_name, 'Unknown') AS role_name").
		Joins("LEFT JOIN role AS r ON u.role_id = r.id").
		Scan(&user).Error

	if err != nil {
		fmt.Println("Error: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve users"})
		return
	}

	if len(user) == 0 {
		c.JSON(http.StatusOK, UserListResponse{Status: "success", Data: []UserDetail{}})
		return
	}

	c.JSON(http.StatusOK, UserListResponse{Status: "success", Data: formatUsers(user)})
}

// CreateUser handles creating a new user with validation
func CreateUser(c *gin.Context, db *gorm.DB) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid input data", "error": err.Error()})
		return
	}

	// Validate username to ensure it's lowercase
	if !isValidLowercaseUsername(newUser.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Username must be in lowercase and can only contain letters, numbers, and underscores"})
		return
	}

	// Validasi username
	if !isValidUsername(newUser.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Username tidak valid"})
		return
	}
	if !isUniqueUsername(db, newUser.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Username sudah digunakan"})
		return
	}

	// Validasi email
	if !isValidEmail(newUser.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Email tidak valid"})
		return
	}
	if !isUniqueEmail(db, newUser.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Email sudah digunakan"})
		return
	}

	// Validasi role_id
	if !isValidRole(db, newUser.RoleID) {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid role_id"})
		return
	}

	// Hash password
	hashedPassword, err := helpers.HashPassword(newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to hash password"})
		return
	}
	newUser.Password = hashedPassword

	// Simpan pengguna baru ke database
	if err := db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create user", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User created successfully!"})
}

// UpdateUser handles updating user details
func UpdateUser(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var user models.User

	// Check if user exists
	if err := db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "User not found"})
		} else {
			fmt.Println("Error finding user: ", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to find user"})
		}
		return
	}

	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		RoleID   int    `json:"role_id"`
		IsActive bool   `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid input data", "error": err.Error()})
		return
	}

	// Validate username to ensure it's lowercase
	if !isValidLowercaseUsername(input.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Username must be in lowercase and can only contain letters, numbers, and underscores"})
		return
	}

	// Validate username format
	if !isValidUsername(input.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid username"})
		return
	}
	// Check if username is unique (excluding current user)
	if !isUniqueUsernameForUpdate(db, input.Username, user.ID) {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Username already in use"})
		return
	}

	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(input.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid email format"})
		return
	}
	// Check if email is unique (excluding current user)
	if !isUniqueEmailForUpdate(db, input.Email, user.ID) {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Email already in use"})
		return
	}

	// Validate role_id
	if !isValidRole(db, input.RoleID) {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid role_id"})
		return
	}

	// Update user fields
	user.Username = input.Username
	user.Email = input.Email
	user.RoleID = input.RoleID
	user.IsActive = input.IsActive
	user.UpdatedAt = time.Now()

	// Save the updated user
	if err := db.Save(&user).Error; err != nil {
		fmt.Println("Error updating user: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to update user", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User updated successfully"})
}

// DeleteUser handles deleting a user by ID
func DeleteUser(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	var user models.User
	if err := db.First(&user, "id = ?", id).Error; err != nil {
		handleError(c, err, "Gagal mendapatkan data user")
		return
	}

	if user.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "User masih aktif, tidak dapat dihapus"})
		return
	}

	// Cek apakah user memiliki data employee atau attendance terkait
	var count int64
	if err := db.Model(&models.Employee{}).Where("user_id = ?", user.ID).Or("user_id = ?", user.ID).Count(&count).Error; err != nil {
		handleError(c, err, "Gagal memeriksa data terkait user")
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "User memiliki data karyawan atau absensi, tidak dapat dihapus"})
		return
	}

	// Periksa role_id yang valid sebelum menghapus user
	if !isValidRole(db, user.RoleID) {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Role terkait user tidak valid"})
		return
	}

	// Hapus user jika semua pengecekan berhasil
	if err := db.Delete(&user).Error; err != nil {
		handleError(c, err, "Gagal menghapus user")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User berhasil dihapus"})
}
// formatUsers formats users with role names
func formatUsers(users []struct {
	models.User
	RoleName string `gorm:"column:role_name"`
}) []UserDetail {
	var formattedUsers []UserDetail
	for _, user := range users {
		formattedUsers = append(formattedUsers, UserDetail{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			RoleID:   user.RoleID,
			RoleName: user.RoleName,
			IsActive: user.IsActive,
		})
	}
	return formattedUsers
}

// isValidUsername checks if username has no spaces
func isValidUsername(username string) bool {
	const usernamePattern = `^\S+$`
	return regexp.MustCompile(usernamePattern).MatchString(username)
}

// isUniqueUsername checks if username is unique in the database
func isUniqueUsername(db *gorm.DB, username string) bool {
	var count int64
	db.Model(&models.User{}).Where("username = ?", username).Count(&count)
	return count == 0
}

// isValidEmail validates the email format
func isValidEmail(email string) bool {
	const emailPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(emailPattern).MatchString(email)
}

// isUniqueEmail checks if email is unique in the database
func isUniqueEmail(db *gorm.DB, email string) bool {
	var count int64
	db.Model(&models.User{}).Where("email = ?", email).Count(&count)
	return count == 0
}

// CheckUsername handles checking if the username is already in use
func CheckUsername(c *gin.Context, db *gorm.DB) {
	var request struct {
		Username string `json:"username"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid input"})
		return
	}

	isUnique := isUniqueUsername(db, request.Username)
	c.JSON(http.StatusOK, gin.H{"status": "success", "username_exists": !isUnique})
}

// CheckEmail handles checking if the email is already in use
func CheckEmail(c *gin.Context, db *gorm.DB) {
	var request struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid input"})
		return
	}

	isUnique := isUniqueEmail(db, request.Email)
	c.JSON(http.StatusOK, gin.H{"status": "success", "email_exists": !isUnique})
}

// isValidLowercaseUsername checks if a username is lowercase and matches required criteria
func isValidLowercaseUsername(username string) bool {
	// Regular expression to check for only lowercase letters
	lowercaseRegex := `^[a-z0-9_]+$`
	matched, _ := regexp.MatchString(lowercaseRegex, username)
	return matched
}

// isValidRole checks if the role exists in the database
func isValidRole(db *gorm.DB, roleID int) bool {
	var count int64
	db.Model(&models.Role{}).Where("id = ?", roleID).Count(&count)
	return count > 0
}

// isUniqueUsernameForUpdate checks if the username is unique, excluding the current user
func isUniqueUsernameForUpdate(db *gorm.DB, username string, userID int) bool {
	var existingUser models.User
	if err := db.Where("username = ? AND id != ?", username, userID).First(&existingUser).Error; err != nil {
		return err == gorm.ErrRecordNotFound
	}
	return false
}

// isUniqueEmailForUpdate checks if the email is unique, excluding the current user
func isUniqueEmailForUpdate(db *gorm.DB, email string, userID int) bool {
	var existingUser models.User
	if err := db.Where("email = ? AND id != ?", email, userID).First(&existingUser).Error; err != nil {
		return err == gorm.ErrRecordNotFound
	}
	return false
}

func handleError(c *gin.Context, err error, message string) {
	log.Printf("Error: %v", err)
	c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": message, "error": err.Error()})
}
