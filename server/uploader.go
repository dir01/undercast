package server

import (
	"fmt"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type uploader struct {
	awsID     string
	awsSecret string
	awsToken  string
	awsRegion string
	s3Bucket    string
}

// NewUploader creates an uploader that is configured for specific aws credentials
func NewUploader(id, secret, token, region, bucket string) *uploader {
	u := uploader{awsID: id, awsSecret: secret, awsToken: token, awsRegion: region, s3Bucket: bucket}
	return &u
}

func (u *uploader) UploadFile(filename string) (string, error) {
	creds := credentials.NewStaticCredentials(u.awsID, u.awsSecret, u.awsToken)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      &u.awsRegion,
	}))
	uploader := s3manager.NewUploader(sess)
	f, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("failed to open file %q, %v", filename, err)
	}

	_, key := path.Split(filename)
	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(u.s3Bucket),
		Key:    aws.String(key),
		Body:   f,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file, %v", err)
	}
	return result.Location, nil
}
