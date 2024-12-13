package utils

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ApplicationPort int
	ApplicationURL  string

	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     int
	DBName     string

	BackupDBHost     string
	BackupDBPort     int
	BackupDBUser     string
	BackupDBPassword string
	BackupDBName     string

	S3Region     string
	S3Endpoint   string
	S3AccessKey  string
	S3SecretKey  string
	S3BucketName string

	TelegramBotAPIKey string
	TelegramChatID    string

	CronSchedule string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config := &Config{
		ApplicationPort: getEnvInt("APPLICATION_PORT", 14705),
		ApplicationURL:  os.ExpandEnv(getEnv("APPLICATION_URL", "")),

		DBUser:     getEnv("DB_USER", ""),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBHost:     getEnv("DB_HOST", ""),
		DBPort:     getEnvInt("DB_PORT", 5433),
		DBName:     getEnv("DB_NAME", ""),

		BackupDBHost:     getEnv("BACKUP_DB_HOST", ""),
		BackupDBPort:     getEnvInt("BACKUP_DB_PORT", 5432),
		BackupDBUser:     getEnv("BACKUP_DB_USER", ""),
		BackupDBPassword: getEnv("BACKUP_DB_PASSWORD", ""),
		BackupDBName:     getEnv("BACKUP_DB_NAME", ""),

		S3Region:     getEnv("S3_REGION", ""),
		S3Endpoint:   getEnv("S3_ENDPOINT", ""),
		S3AccessKey:  getEnv("S3_ACCESS_KEY", ""),
		S3SecretKey:  getEnv("S3_SECRET_KEY", ""),
		S3BucketName: getEnv("S3_BUCKET_NAME", ""),

		TelegramBotAPIKey: getEnv("TELEGRAM_BOT_API_KEY", ""),
		TelegramChatID:    getEnv("TELEGRAM_CHAT_ID", ""),

		CronSchedule: getEnv("CRON_SCHEDULE", "0 0 * * *"),
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		} else {
			log.Printf("Invalid integer value for %s: %s, using default %d", key, value, defaultValue)
		}
	}
	return defaultValue
}
