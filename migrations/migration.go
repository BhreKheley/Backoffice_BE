package migrations

import (
	"log"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"absensi-app/database"  // Ganti dengan path yang sesuai ke package database kamu
	"absensi-app/models"    // Ganti dengan path yang sesuai ke package models kamu
)

// Migrate function to run migrations
func Migrate() error {
	// Ambil koneksi database yang sudah diinisialisasi
	db := database.InitDB()

	// Definisikan migrasi tabel
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "20241113001", // Ganti dengan ID yang sesuai, bisa pakai timestamp atau lainnya
			Migrate: func(tx *gorm.DB) error {
				// Melakukan AutoMigrate untuk tabel-tabel
				return tx.AutoMigrate(
					&models.Attendance{},
					&models.Division{},
					&models.Employee{},
					&models.Permission{},
					&models.Position{},
					&models.Role{},
					&models.RolePermission{},
					&models.Status{},
					&models.User{},
				)
			},
			Rollback: func(tx *gorm.DB) error {
				// Rollback tabel jika dibutuhkan
				return tx.Migrator().DropTable(
					"attendance",
					"division",
					"employee",
					"permission",
					"position",
					"role",
					"role_permission",
					"status",
					"user",
				)
			},
		},
	})

	// Jalankan migrasi
	if err := m.Migrate(); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
		return err
	}

	log.Println("Migrations applied successfully!")
	return nil
}
