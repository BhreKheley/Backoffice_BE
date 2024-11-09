package migrations

import (
	"gorm.io/gorm"
	// "absensi-app/models"
)

func Migrate(db *gorm.DB) {
	// Migrasi untuk Role, Permission, dan RolePermission (relasi many-to-many)
	// db.AutoMigrate(&models.Role{}, &models.Permission{}, &models.RolePermission{})
}
