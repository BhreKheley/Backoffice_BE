package handlers

import (
	"absensi-app/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// GetRoles - Mendapatkan semua role
func GetRoles(c *gin.Context, db *sqlx.DB) {
    var roles []models.Role
    query := `SELECT id, role_name, code FROM role`
    
    log.Println("Executing query to fetch roles: ", query) // Tambahkan ini

    err := db.Select(&roles, query)
    if err != nil {
        log.Printf("Error fetching roles: %v\n", err) // Tambahkan ini
        c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to fetch roles"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"status": "success", "data": roles})
}


// CreateRole - Membuat role baru
func CreateRole(c *gin.Context, db *sqlx.DB) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request"})
		return
	}

	query := `INSERT INTO role (role_name, code) VALUES ($1, $2) RETURNING id`
	err := db.QueryRow(query, role.RoleName, role.Code).Scan(&role.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": role})
}

// GetPermissions - Mendapatkan semua permissions
func GetPermissions(c *gin.Context, db *sqlx.DB) {
	var permissions []models.Permission
	query := `SELECT id, permission_name, code FROM permission`

	err := db.Select(&permissions, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to fetch permissions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": permissions})
}

// CreatePermission - Membuat permission baru
func CreatePermission(c *gin.Context, db *sqlx.DB) {
	var permission models.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request"})
		return
	}

	query := `INSERT INTO permission (permission_name, code) VALUES ($1, $2) RETURNING id`
	err := db.QueryRow(query, permission.PermissionName, permission.Code).Scan(&permission.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create permission"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": permission})
}

// AssignPermissionToRole - Menambahkan permission ke role
func AssignPermissionToRole(c *gin.Context, db *sqlx.DB) {
	var rolePermission models.RolePermission
	if err := c.ShouldBindJSON(&rolePermission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request"})
		return
	}

	query := `INSERT INTO role_permission (role_id, permission_id) VALUES ($1, $2)`
	_, err := db.Exec(query, rolePermission.RoleID, rolePermission.PermissionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to assign permission to role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Permission assigned to role"})
}

// GetPermissionsByRole - Mendapatkan permissions berdasarkan role
func GetPermissionsByRole(c *gin.Context, db *sqlx.DB) {
	roleID := c.Param("roleID")
	var permissions []models.Permission

	query := `
		SELECT p.id, p.permission_name, p.code
		FROM permission p
		INNER JOIN role_permission rp ON rp.permission_id = p.id
		WHERE rp.role_id = $1`

	err := db.Select(&permissions, query, roleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to fetch permissions for the role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": permissions})
}

// RemovePermissionFromRole - Menghapus permission dari role
func RemovePermissionFromRole(c *gin.Context, db *sqlx.DB) {
	var rolePermission models.RolePermission
	if err := c.ShouldBindJSON(&rolePermission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request"})
		return
	}

	query := `DELETE FROM role_permission WHERE role_id = $1 AND permission_id = $2`
	_, err := db.Exec(query, rolePermission.RoleID, rolePermission.PermissionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to remove permission from role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Permission removed from role"})
}
