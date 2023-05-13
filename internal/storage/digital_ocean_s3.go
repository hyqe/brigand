package storage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
)

func newS3Session(region, s3_endpoint, accessKey, secretKey string) (*session.Session, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:    aws.String(s3_endpoint),
	}))

	return sess, nil
}

type FileUploader func(file io.Reader, filename string) error

func s3FileUploader(sess *session.Session, file io.Reader, filename string) error {
	_, err := s3manager.NewUploader(sess).Upload(&s3manager.UploadInput{
		Key:  aws.String(filename),
		Body: file,
	})

	return err
}

func NewS3FileUploader(sess *session.Session) FileUploader {
	return func(file io.Reader, filename string) error {
		return s3FileUploader(sess, file, filename)
	}
}

type FileDownloader func(file io.Writer, filename string) error

func s3FileDownloader(sess *session.Session, file io.Writer, filename string) error {
	r, err := s3.New(sess).GetObject(&s3.GetObjectInput{
		Key: aws.String(filename),
	})
	if err != nil {
		return err
	}

	io.Copy(file, r.Body)

	return err
}

func NewS3FileDownloader(s3Sess *session.Session) FileDownloader {
	return func(file io.Writer, filename string) error {
		return s3FileDownloader(s3Sess, file, filename)
	}
}