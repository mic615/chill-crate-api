package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/mic615/chill-crate-api/internal/config"
)

var Client *s3.Client

func Connect(cfg *config.Config) {
	awscfg, err := awsconfig.LoadDefaultConfig(
		context.Background(),
		awsconfig.WithRegion(cfg.StorageRegion),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.StorageAccessKey, cfg.StorageSecretKey, "")),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(awscfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.StorageEndpoint)
		o.UsePathStyle = true
	})
	client.ListBuckets(context.Background(), nil)
	// reachability probe — fail fast at boot.
	if _, err := client.ListBuckets(context.Background(), &s3.ListBucketsInput{}); err != nil {
		log.Fatalf("failed to reach Storage endpoint %s: %v", cfg.StorageEndpoint, err)
	}
	Client = client
}

func CreateBucket(name string) error {
	_, err := Client.CreateBucket(context.Background(), &s3.CreateBucketInput{Bucket: aws.String(name)})
	if err != nil {
		return fmt.Errorf("create bucket %s: %w", name, err)
	}
	return err
}

func UploadObject(bucketName string, objectKey string, fileName string, file multipart.File) error {
	defer file.Close()
	_, err := Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("Couldn't upload file %v to %v:%v. Here's why: %v\n", fileName, bucketName, objectKey, err)
	}
	return err
}

func DownloadObject(bucketName string, objectKey string) (io.ReadCloser, error) {
	result, err := Client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
	})
	if err != nil {
		return nil, fmt.Errorf("Couldn't get object %v:%v. Here's why: %v\n", bucketName, objectKey, err)

	}
	return result.Body, err

}
