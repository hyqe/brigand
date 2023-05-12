package storage

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"

	"os"
)

const (
	BUCKET_NAME = "brigand-storage"
	REGION      = "US"
)

func getEnv(env_name string) (string, error) {
	env, ok := os.LookupEnv(env_name)
	if !ok {
		return "", fmt.Errorf("There is an error! No %s!!!!", env_name)
	}

	return env, nil
}

func newS3Session() (*session.Session, error) {
	envs := map[string]string{
		"SECRET_KEY": "",
		"ACCESS_KEY": "",
		"ENDPOINT":   "",
	}

	for name := range envs {
		env, err := getEnv(name)
		if err != nil {
			return nil, err
		}
		envs[name] = env

	}

	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(REGION),
		Credentials: credentials.NewStaticCredentials(envs["ACCESS_KEY"], envs["SECRET_KEY"], ""),
		Endpoint:    aws.String(envs["ENDPOINT"]),
	}))

	return sess, nil
}

type FileUploader func(file io.Reader, filename string) error

func s3FileUploader(sess *session.Session, file io.Reader, filename string) error {
	_, err := s3manager.NewUploader(sess).Upload(&s3manager.UploadInput{
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(filename),
		Body:   file,
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
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(filename),
	})
	if err != nil {
		return err
	}

	io.Copy(file, r.Body)

	return err
}

func NewS3FileDownloader(sess *session.Session) FileDownloader {
	return func(file io.Writer, filename string) error {
		return s3FileDownloader(sess, file, filename)
	}
}