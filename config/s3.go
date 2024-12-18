package config

import (
	"backup/utils"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewS3Client(cfg *utils.Config) *s3.Client {
	s3Config := aws.Config{
		Region:       cfg.S3Region,
		Credentials:  credentials.NewStaticCredentialsProvider(cfg.S3AccessKey, cfg.S3SecretKey, ""),
		BaseEndpoint: &cfg.S3Endpoint,
	}

	client := s3.NewFromConfig(s3Config)
	log.Println("✅ S3 client created successfully")

	return client
}
