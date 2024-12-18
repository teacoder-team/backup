package main

import (
	"backup/config"
	"backup/services"
	"backup/utils"
	"log"
)

func main() {
	cfg, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	config.ConnectDatabase(cfg)

	backupDB, err := config.ConnectBackupDatabase(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to backup database: %v", err)
	}

	s3Client := config.NewS3Client(cfg)

	cronService := services.NewCronService(cfg, config.DB, backupDB, s3Client)
	cronService.Start()

	log.Println("✅ Backup service started successfully!")

	select {}
}
