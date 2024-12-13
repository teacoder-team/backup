package models

import (
	"time"
)

type BackupStatus string

const (
	BackupStatusSuccess    BackupStatus = "success"
	BackupStatusFailed     BackupStatus = "failed"
	BackupStatusInProgress BackupStatus = "in_progress"
)

type Backup struct {
	ID         uint          `gorm:"primaryKey;column:id"`
	CreatedAt  time.Time     `gorm:"column:created_at"`
	BackupTime time.Time     `gorm:"column:backup_time"`
	Status     BackupStatus  `gorm:"type:text;column:status;default:in_progress"`
	S3Path     string        `gorm:"size:255;column:s3_path"`
	BackupSize int64         `gorm:"column:backup_size"`
	Duration   time.Duration `gorm:"column:duration"`
}
