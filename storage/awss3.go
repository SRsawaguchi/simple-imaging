package storage

import (
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var prefix = "result/"

type AwsS3 struct {
	Downloader
	Uploader
	Region          string
	Bucket          string
	session         *session.Session
	tempFileManager TempFileManager
}

func NewAwsS3(region, bucket string) *AwsS3 {
	return &AwsS3{
		Region: region,
		Bucket: bucket,
	}
}

func NewAwsS3FromEnvironment() *AwsS3 {
	return &AwsS3{
		Region: os.Getenv("SRIMAGE_S3_REGION"),
		Bucket: os.Getenv("SRIMAGE_S3_BUCKET"),
	}
}

func (s *AwsS3) connect() error {
	if s.session != nil {
		return nil
	}

	session, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(s.Region)},
	})
	if err != nil {
		return err
	}
	s.session = session

	return nil
}

func (s AwsS3) GenerateUploadKey(fileName string) string {
	return "grayscale/" + fileName
}

func (s AwsS3) Upload(key string, src io.Reader) error {
	if err := s.connect(); err != nil {
		return err
	}
	uploader := s3manager.NewUploader(s.session)

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
		Body:   src,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s AwsS3) Download(key string, dest io.WriterAt) error {
	if err := s.connect(); err != nil {
		return err
	}

	downloader := s3manager.NewDownloader(s.session)
	_, err := downloader.Download(dest,
		&s3.GetObjectInput{
			Bucket: aws.String(s.Bucket),
			Key:    aws.String(key),
		},
	)
	if err != nil {
		return err
	}

	return nil
}
