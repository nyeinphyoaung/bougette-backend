package services

import (
	"bougette-backend/configs"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type UploadService struct {
	s3Client   *s3.Client
	bucketName string
}

func NewUploadService() *UploadService {
	awsConfig, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(configs.Envs.AWS_REGION),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			configs.Envs.AWS_ACCESS_KEY_ID,
			configs.Envs.AWS_SECRET_ACCESS_KEY,
			"",
		)),
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to load AWS config: %v", err))
	}

	return &UploadService{
		s3Client:   s3.NewFromConfig(awsConfig),
		bucketName: configs.Envs.AWS_BUCKET_NAME,
	}
}

func (s *UploadService) GeneratePresignedUploadURL(filename string) (string, string, error) {
	presignClient := s3.NewPresignClient(s.s3Client)
	key := fmt.Sprintf("%s-%s", uuid.New().String(), filename)

	request, err := presignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(3600 * time.Second)
	})

	if err != nil {
		return "", "", fmt.Errorf("failed to generate presigned upload URL: %w", err)
	}

	return request.URL, key, nil
}

func (s *UploadService) GeneratePresignedDownloadURL(key string) (string, error) {
	presignClient := s3.NewPresignClient(s.s3Client)

	request, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(3600 * time.Second)
	})

	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return request.URL, nil
}

func (s *UploadService) DeleteFile(key string) error {
	_, err := s.s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	return nil
}
