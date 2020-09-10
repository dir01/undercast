package undercast

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"os"
	"path"
)

type s3Storage struct {
	s3Config   *aws.Config
	s3Client   *s3.S3
	keyPrefix  string
	bucketName string
}

func (storage *s3Storage) Store(ctx context.Context, filepath, filename string) (url string, err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	return storage.StoreData(ctx, file, filename)
}

func (storage *s3Storage) StoreData(ctx context.Context, data io.ReadSeeker, filename string) (url string, err error) {
	key := path.Join(storage.keyPrefix, filename)
	input := &s3.PutObjectInput{
		Body:   data,
		Bucket: &storage.bucketName,
		Key:    &key,
	}
	if _, err := storage.s3Client.PutObjectWithContext(ctx, input); err != nil {
		return "", err
	}
	return formatS3Url(*storage.s3Config.Region, storage.bucketName, key), nil

}

func formatS3Url(region, bucket, key string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, key)
}
