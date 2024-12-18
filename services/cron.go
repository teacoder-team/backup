package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"backup/models"
	"backup/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type CronService struct {
	DB          *gorm.DB
	BackupDB    *gorm.DB
	S3Client    *s3.Client
	Config      *utils.Config
	TelegramBot *tgbotapi.BotAPI
}

func NewCronService(cfg *utils.Config, db *gorm.DB, backupDB *gorm.DB, s3Client *s3.Client) *CronService {
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBotAPIKey)
	if err != nil {
		log.Fatalf("❌ Failed to initialize Telegram bot: %v", err)
	}

	return &CronService{
		DB:          db,
		BackupDB:    backupDB,
		S3Client:    s3Client,
		Config:      cfg,
		TelegramBot: bot,
	}
}

func (s *CronService) Start() {
	c := cron.New()

	_, err := c.AddFunc(s.Config.CronSchedule, func() {
		err := s.PerformBackup()
		if err != nil {
			log.Printf("❌ Backup failed: %v", err)
		} else {
			log.Println("✅ Backup completed successfully")
			s.SendTelegramNotification()
		}
	})
	if err != nil {
		log.Fatalf("❌ Failed to add cron job: %v", err)
	}

	log.Println("✅ Cron job started successfully")
	c.Start()
}

func (s *CronService) PerformBackup() error {
	var tables []string
	err := s.BackupDB.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'").Scan(&tables).Error
	if err != nil {
		return fmt.Errorf("failed to fetch tables from backup database: %v", err)
	}

	backupTime := time.Now()
	s3Path := fmt.Sprintf("%s", backupTime.Format("2006-01-02"))

	_, err = s.S3Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:       aws.String(s.Config.S3BucketName),
		Key:          aws.String(s3Path + "/"),
		Body:         bytes.NewReader([]byte{}),
		StorageClass: "GLACIER",
	})
	if err != nil {
		return fmt.Errorf("failed to create backup folder in S3: %v", err)
	}

	var backupSize int64
	for _, table := range tables {
		var data []map[string]interface{}
		err := s.BackupDB.Table(table).Find(&data).Error
		if err != nil {
			log.Printf("Failed to fetch data for table %s: %v", table, err)
			continue
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Printf("Failed to marshal data for table %s: %v", table, err)
			continue
		}

		key := fmt.Sprintf("%s/%s.json", s3Path, table)
		_, err = s.S3Client.PutObject(context.Background(), &s3.PutObjectInput{
			Bucket: aws.String(s.Config.S3BucketName),
			Key:    aws.String(key),
			Body:   bytes.NewReader(jsonData),
		})
		if err != nil {
			log.Printf("Failed to upload table %s to S3: %v", table, err)
			continue
		}

		backupSize += int64(len(jsonData))
	}

	backup := models.Backup{
		BackupTime: backupTime,
		Status:     models.BackupStatusSuccess,
		S3Path:     s3Path,
		BackupSize: backupSize,
		Duration:   time.Since(backupTime),
	}

	err = s.DB.Create(&backup).Error
	if err != nil {
		return fmt.Errorf("failed to record backup in database: %v", err)
	}

	log.Printf("Backup completed successfully: %s", s3Path)
	return nil
}

func (s *CronService) SendTelegramNotification() {
	messageText := fmt.Sprintf(
		"✅ Выполнена резервная копия базы данных!",
	)

	msg := tgbotapi.NewMessageToChannel(s.Config.TelegramChatID, messageText)

	_, err := s.TelegramBot.Send(msg)
	if err != nil {
		log.Printf("❌ Failed to send Telegram message: %v", err)
		return
	}

	log.Println("✅ Telegram notification sent")
}
