package handlers

import (
	"absensi-app/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetRoles - Mendapatkan semua role
func GetRoles(c *gin.Context, db *gorm.DB) {
	var roles []models.Role
	log.Println("Executing query to fetch roles")

	if err := db.Find(&roles).Error; err != nil {
		log.Printf("Error fetching roles: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to fetch roles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": roles})
}

// CreateRole - Membuat role baru
func CreateRole(c *gin.Context, db *gorm.DB) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request"})
		return
	}

	if err := db.Create(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": role})
}

// GetPermissions - Mendapatkan semua permissions
func GetPermissions(c *gin.Context, db *gorm.DB) {
	var permissions []models.Permission

	if err := db.Find(&permissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to fetch permissions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": permissions})
}

// CreatePermission - Membuat permission baru
func CreatePermission(c *gin.Context, db *gorm.DB) {
	var permission models.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request"})
		return
	}

	if err := db.Create(&permission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create permission"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": permission})
}

// AssignPermissionToRole - Menambahkan permission ke role
func AssignPermissionToRole(c *gin.Context, db *gorm.DB) {
	var rolePermission models.RolePermission
	if err := c.ShouldBindJSON(&rolePermission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request"})
		return
	}

	if err := db.Create(&rolePermission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to assign permission to role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Permission assigned to role"})
}

// GetPermissionsByRole - Mendapatkan permissions berdasarkan role
func GetPermissionsByRole(c *gin.Context, db *gorm.DB) {
	roleID := c.Param("roleID")
	var permissions []models.Permission

	if err := db.Joins("INNER JOIN role_permission rp ON rp.permission_id = permissions.id").
		Where("rp.role_id = ?", roleID).
		Find(&permissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to fetch permissions for the role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": permissions})
}

// RemovePermissionFromRole - Menghapus permission dari role
func RemovePermissionFromRole(c *gin.Context, db *gorm.DB) {
	var rolePermission models.RolePermission
	if err := c.ShouldBindJSON(&rolePermission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request"})
		return
	}

	if err := db.Where("role_id = ? AND permission_id = ?", rolePermission.RoleID, rolePermission.PermissionID).
		Delete(&models.RolePermission{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to remove permission from role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Permission removed from role"})
}
