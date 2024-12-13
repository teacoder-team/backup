package config

import (
	"example/m/models"
	"example/m/utils"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var BackupDB *gorm.DB

func ConnectDatabase(cfg *utils.Config) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to establish database connection: %v", err)
	}

	log.Println("✅ Database connection established successfully")

	err = DB.AutoMigrate(&models.Backup{})
	if err != nil {
		log.Fatalf("❌ Error during migration: %v", err)
	}

	log.Println("✅ Database migrated successfully")
}

func ConnectBackupDatabase(cfg *utils.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		cfg.BackupDBHost,
		cfg.BackupDBUser,
		cfg.BackupDBPassword,
		cfg.BackupDBName,
		cfg.BackupDBPort,
	)

	var err error
	BackupDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("❌ Failed to establish backup database connection: %v", err)
		return nil, err
	}

	log.Println("✅ Backup database connection established successfully")
	return BackupDB, nil
}
