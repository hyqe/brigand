package handlers_test

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/hyqe/brigand/internal/handlers"
	"github.com/hyqe/brigand/internal/storage"

	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"

	"testing"
)

func getEnv(env_name string) (string, error) {
	env, ok := os.LookupEnv(env_name)
	if !ok {
		return "", fmt.Errorf("There is an error! No %s!!!!", env_name)
	}

	return env, nil
}

func getSomeEnvs() (map[string]string, error) {
	envs := map[string]string{
		"REGION":      "",
		"S3_ENDPOINT": "",
		"ACCESS_KEY":  "",
		"SECRET_KEY":  "",
	}

	for name := range envs {
		env, err := getEnv(name)
		if err != nil {
			return map[string]string{}, err
		}
		envs[name] = env
	}
	return envs, nil
}

func testSession() (*session.Session, error) {
	e, err := getSomeEnvs()
	if err != nil {
		return nil, err
	}

	s3sess, err := storage.NewS3Session(e["REGION"], e["S3_ENDPOINT"], e["ACCESS_KEY"], e["SECRET_KEY"])
	if err != nil {
		return nil, err
	}

	return s3sess, nil
}

func createFile(s3sess *session.Session, filename string, file io.Reader, bucket string) error {
	err := storage.NewS3FileUploader(s3sess, bucket)(file, filename)

	return err
}

func deleteImage(filename string, s3Sess *session.Session) error {
	bucket, err := getEnv("BUCKET")
	if err != nil {
		return err
	}

	s3Client := s3.New(s3Sess)
	_, err = s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})

	return err
}

func Test_GetFile_happy_path(t *testing.T) {
	MONGO, ok := os.LookupEnv("MONGO")
	if !ok {
		t.Skipf("Missing Env: MONGO")
	}

	BUCKET, ok := os.LookupEnv("BUCKET")
	require.True(t, ok)
	// Insert file into s3 and mongodb

	// Get a Mongo Client
	ctx := context.Background()
	mongoClient, err := storage.NewMongoClient(ctx, MONGO)
	require.NoError(t, err)
	defer mongoClient.Disconnect(ctx)

	// Get Metadata Client
	mdClient := storage.NewMongoMetadataClient(mongoClient)

	// INSERT into a record into MongoDb
	md := storage.NewMetadata("Ryujin-Breaker")
	require.NoError(t, mdClient.Create(ctx, md))

	// Create a S3/D.O. Spaces Session
	s3sess, err := testSession()
	if err != nil {
		require.NoError(t, err)
	}

	// INSERT a File into S3/Digital Ocean Spaces
	file := []byte("string")
	reader := bytes.NewReader(file)
	require.NoError(t, createFile(s3sess, md.Id, reader, BUCKET))

	// Get File from S3/D.O. Spaces
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/files/{fileId}", nil)
	r = mux.SetURLVars(r, map[string]string{"fileId": md.Id})
	handlers.NewGetFileById(mdClient, storage.NewS3FileDownloader(s3sess, BUCKET)).ServeHTTP(w, r)

	// Check the response
	if w.Result().StatusCode < 200 || w.Result().StatusCode > 299 {
		require.True(t, false)
	}

	// Compare them, They should be equal
	newfile, err := ioutil.ReadAll(w.Body)
	require.True(t, bytes.Equal(file, newfile))

	// Clean Up
	// // Delete from mongodb
	require.NoError(t, mdClient.DeleteById(ctx, md.Id))
	// // Delete from S3/spaces
	require.NoError(t, deleteImage(md.Id, s3sess))
}

// Test_GetFile_happy_path(t *testing.T) {
func Test_GetFileById_bad_file_id(t *testing.T) {
	MONGO, ok := os.LookupEnv("MONGO")
	require.True(t, ok)

	BUCKET, ok := os.LookupEnv("BUCKET")
	require.True(t, ok)

	// Get a Mongo Client
	ctx := context.Background()
	mongoClient, err := storage.NewMongoClient(ctx, MONGO)
	require.NoError(t, err)
	defer mongoClient.Disconnect(ctx)

	// Get Metadata Client
	mdClient := storage.NewMongoMetadataClient(mongoClient)

	// // INSERT into a record into MongoDb
	md := storage.NewMetadata("Ryujin-Breaker")
	require.NoError(t, mdClient.Create(ctx, md))

	// // CLEANUP: Delete from mongodb
	defer require.NoError(t, mdClient.DeleteById(ctx, md.Id))

	// Create a S3/D.O. Spaces Session
	s3sess, err := testSession()
	if err != nil {
		require.NoError(t, err)
	}

	// INSERT a File into S3/Digital Ocean Spaces
	file := []byte("string")
	reader := bytes.NewReader(file)
	require.NoError(t, createFile(s3sess, md.Id, reader, BUCKET))

	// // CLEANUP: Delete from S3/spaces
	defer require.NoError(t, deleteImage(md.Id, s3sess))

	// Get File from S3/D.O. Spaces
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/files/{fileId}", nil)
	// Use a bad ID
	r = mux.SetURLVars(r, map[string]string{"fileId": uuid.New().String()})
	handlers.NewGetFileById(mdClient, storage.NewS3FileDownloader(s3sess, BUCKET)).ServeHTTP(w, r)

	// Check the response
	require.Equal(t, 404, w.Result().StatusCode)
}
