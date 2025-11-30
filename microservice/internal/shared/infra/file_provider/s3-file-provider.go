package file_service

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"tech_challenge/internal/shared/config/env"
)

type S3FileProvider struct {
	client     *s3.Client
	bucketName string
}

func NewS3FileProvider(bucketName string) *S3FileProvider {
	cfgEnv := env.GetConfig()

	if bucketName == "" {
		bucketName = cfgEnv.AWS.S3.BucketName
	}

	customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		if service == s3.ServiceID {
			return aws.Endpoint{
				URL:               cfgEnv.AWS.S3.Endpoint,
				SigningRegion:     region,
				HostnameImmutable: true,
			}, nil
		}

		return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
	})

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(cfgEnv.AWS.Region),
		config.WithEndpointResolver(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfgEnv.AWS.AccessKeyID,
			cfgEnv.AWS.SecretAccessKey,
			"",
		)),
	)

	if err != nil {
		return nil
	}

	return &S3FileProvider{
		client:     s3.NewFromConfig(cfg),
		bucketName: bucketName,
	}
}

func (s *S3FileProvider) UploadFile(fileName string, fileContent []byte) error {
	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fileName),
		Body:   bytes.NewReader(fileContent),
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *S3FileProvider) DeleteFile(fileName string) error {
	_, err := s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fileName),
	})

	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

func (s *S3FileProvider) GetPresignedURL(fileName string) (string, error) {
	presignClient := s3.NewPresignClient(s.client)

	presignedRequest, err := presignClient.PresignGetObject(
		context.TODO(),
		&s3.GetObjectInput{
			Bucket: aws.String(s.bucketName),
			Key:    aws.String(fileName),
		},
		func(o *s3.PresignOptions) {
			o.Expires = 15 * time.Minute
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to get presigned URL: %w", err)
	}

	return presignedRequest.URL, nil
}
