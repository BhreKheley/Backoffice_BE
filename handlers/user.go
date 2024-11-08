package handlers

import (
	"absensi-app/helpers"
	"absensi-app/models"
	"fmt"
	"net/http"
	"regexp"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
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

// GetUserByID handles fetching user by ID
func GetUserByID(c *gin.Context, db *sqlx.DB) {
	id := c.Param("id")

	var user struct {
		models.User
		RoleName string `db:"role_name"`
	}

	err := db.QueryRow(`
		SELECT j1.id, j1.username, j1.email, j1.role_id, j1.is_active, j2.role_name 
		FROM "user" j1
		LEFT JOIN role j2 ON j1.role_id = j2.id 
		WHERE j1.id = $1`, id).
		Scan(&user.ID, &user.Username, &user.Email, &user.RoleID, &user.IsActive, &user.RoleName)

	if err != nil {
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
func GetAllUsers(c *gin.Context, db *sqlx.DB) {
	var users []struct {
		models.User
		RoleName string `db:"role_name"`
	}

	err := db.Select(&users, `
		SELECT j1.id, j1.username, j1.email, j1.role_id, j1.is_active, j2.role_name 
		FROM "user" j1
		LEFT JOIN role j2 ON j1.role_id = j2.id`)

	if err != nil {
		fmt.Println("Error: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve users"})
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusOK, UserListResponse{Status: "success", Data: []UserDetail{}})
		return
	}

	c.JSON(http.StatusOK, UserListResponse{Status: "success", Data: formatUsers(users)})
}

// CreateUser handles creating a new user with validation
func CreateUser(c *gin.Context, db *sqlx.DB) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid input data", "error": err.Error()})
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

	// Hash password
	hashedPassword, err := helpers.HashPassword(newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to hash password"})
		return
	}

	// Simpan pengguna baru ke database
	_, err = db.Exec(`INSERT INTO "user" (username, email, password, role_id, is_active) VALUES ($1, $2, $3, $4, $5)`,
		newUser.Username, newUser.Email, hashedPassword, newUser.RoleID, newUser.IsActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create user", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User created successfully!"})
}


// UpdateUser handles updating an existing user by ID
func UpdateUser(c *gin.Context, db *sqlx.DB) {
	id := c.Param("id")
	var updatedUser models.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid input data", "error": err.Error()})
		return
	}

	_, err := db.Exec(`UPDATE "user" SET username = $1, email = $2, role_id = $3, is_active = $4 WHERE id = $5`,
		updatedUser.Username, updatedUser.Email, updatedUser.RoleID, updatedUser.IsActive, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to update user", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User updated successfully!"})
}

// DeleteUser handles deleting a user by ID
func DeleteUser(c *gin.Context, db *sqlx.DB) {
	id := c.Param("id")

	// Cek apakah user ada di database
	var user models.User
	err := db.Get(&user, `SELECT * FROM "user" WHERE id = $1`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "User tidak ditemukan"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mendapatkan data user", "error": err.Error()})
		}
		return
	}

	// Cek apakah user memiliki data employee terkait
	var employeeCount int
	err = db.Get(&employeeCount, `SELECT COUNT(*) FROM employee WHERE user_id = $1`, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal memeriksa data karyawan", "error": err.Error()})
		return
	}

	// Cek apakah user memiliki data attendance terkait
	var attendanceCount int
	err = db.Get(&attendanceCount, `SELECT COUNT(*) FROM attendance WHERE user_id = $1`, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal memeriksa data absensi", "error": err.Error()})
		return
	}

	// Jika user memiliki data di employee atau attendance, batalkan penghapusan
	if employeeCount > 0 || attendanceCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "User memiliki data karyawan atau absensi, tidak dapat dihapus"})
		return
	}

	// Jika user belum memiliki data employee dan attendance, lanjutkan hanya jika user tidak aktif
	if user.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "User masih aktif, tidak dapat dihapus"})
		return
	}

	// Hapus user jika tidak memiliki data terkait dan statusnya tidak aktif
	_, err = db.Exec(`DELETE FROM "user" WHERE id = $1`, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menghapus user", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User berhasil dihapus"})
}

// formatUsers formats users with role names
func formatUsers(users []struct {
	models.User
	RoleName string `db:"role_name"`
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
func isUniqueUsername(db *sqlx.DB, username string) bool {
	var count int
	err := db.Get(&count, `SELECT COUNT(*) FROM "user" WHERE username = $1`, username)
	return err == nil && count == 0
}

// isValidEmail validates the email format
func isValidEmail(email string) bool {
	const emailPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(emailPattern).MatchString(email)
}

// isUniqueEmail checks if email is unique in the database
func isUniqueEmail(db *sqlx.DB, email string) bool {
	var count int
	err := db.Get(&count, `SELECT COUNT(*) FROM "user" WHERE email = $1`, email)
	return err == nil && count == 0
}

// CheckUsername handles checking if the username is already in use
func CheckUsername(c *gin.Context, db *sqlx.DB) {
    var request struct {
        Username string `json:"username"`
    }

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid input data", "error": err.Error()})
        return
    }

    var usernameCount int
    if err := db.Get(&usernameCount, `SELECT COUNT(*) FROM "user" WHERE username = $1`, request.Username); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to check username", "error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":          "success",
        "username_exists": usernameCount > 0,
    })
}

// CheckEmail handles checking if the email is already in use
func CheckEmail(c *gin.Context, db *sqlx.DB) {
    var request struct {
        Email string `json:"email"`
    }

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid input data", "error": err.Error()})
        return
    }

    var emailCount int
    if err := db.Get(&emailCount, `SELECT COUNT(*) FROM "user" WHERE email = $1`, request.Email); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to check email", "error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":       "success",
        "email_exists": emailCount > 0,
    })
}


