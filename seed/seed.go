package seed

import (
	"gorm.io/gorm"
	// "absensi-app/models"
)

func SeedRolesAndPermissions(db *gorm.DB) {
	// // Buat roles
	// superAdminRole := models.Role{RoleName: "Super Admin", Code: "super_admin"}
	// accountingRole := models.Role{RoleName: "Accounting", Code: "accounting"}

	// // Buat permissions
	// viewAllAttendance := models.Permission{PermissionName: "View All Attendance", Code: "view_all_attendance"}
	// viewOwnAttendance := models.Permission{PermissionName: "View Own Attendance", Code: "view_own_attendance"}
	// viewEmployeeData := models.Permission{PermissionName: "View Employee Data", Code: "view_employee_data"}

	// // Simpan roles dan permissions ke dalam database
	// db.Create(&superAdminRole)
	// db.Create(&accountingRole)
	// db.Create(&viewAllAttendance)
	// db.Create(&viewOwnAttendance)
	// db.Create(&viewEmployeeData)

	// // Assign permissions ke roles
	// db.Model(&superAdminRole).Association("Permissions").Append(&viewAllAttendance, &viewEmployeeData)
	// db.Model(&accountingRole).Association("Permissions").Append(&viewOwnAttendance)
}
