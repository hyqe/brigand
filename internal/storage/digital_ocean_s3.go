package storage

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	// "github.com/aws/aws-sdk-go/service/s3"
	// "github.com/aws/aws-sdk-go/service/s3/s3manager"

	// "io"
	"os"
)

const (
	BUCKET_NAME = "brigand-storage"
	REGION      = "US"
)

// README
// Digital Ocean is making a weird abstraction over only the gui side of aws s3 as far as I can tell
// You interact with (digital ocean spaces s3) programmatically the same way you would a normal aws s3
// Digital Ocean will provide you with the necessary credentials (which you can supply at runtime as Envs)

// Go to DigitialOcean.com -> Login -> API(on the left bottom side of the website) -> Click "Spaces Keys"
// // From here you can generate a new key or use the current ones they've given you

// You NEED:
// // 1) Access Key; This is always visible on the website. You can reaccess this whenever you need it.
// // 2) Secret Key; This will disapear after some time or after you leave the page or something.
// // // its currently labeled just as "Secret"... Make sure to save it.

// Aws lets you generate a crentials object at runtime. You just need to provide it access_key, secret_key

// BUCKET_NAME:
// The bucket_name is the name of the "Digital Ocean Spaces" Space -- i.e. brigand-storage

// ENDPOINT:
// In order to query the s3 you need to set the endpoint.
// Go to DigitialOcean.com -> Login -> Spaces (leftside navbar) -> Click on storage name -> Origin-Endpoint (copypasta)
// Remove the "<your_storage_spaces_name>." from the URL, the aws session object will reconfigure the endpoint
// // to have your bucket_name put back into the endpoint url. Will not work otherwise, you will have a dumb endpoint

func newS3Session(endpoint string) (*session.Session, error) {
	SECRET_KEY, ok := os.LookupEnv("SECRET_KEY")
	if !ok {
		return nil, fmt.Errorf("There is an error! No SECRET KEY!!!!")
	}

	ACCESS_KEY, ok := os.LookupEnv("ACCESS_KEY")
	if !ok {
		return nil, fmt.Errorf("There is an error! No ACCESS_KEY!!!!")
	}

	// ENDPOINT, ok := os.LookupEnv("ENDPOINT")
	// if !ok {
	// 	return nil, fmt.Errorf("There is an error! No ENDPOINT!!!!")
	// }

	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(REGION),
		Credentials: credentials.NewStaticCredentials(ACCESS_KEY, SECRET_KEY, ""),
		Endpoint:    aws.String(endpoint),
	}))

	return sess, nil

}

// func s3UploadFile(sess *session.Session, file *io.Reader, filename string) error {
// _, err := s3manager.NewUploader(sess).Upload(&s3manager.UploadInput{
// Bucket: aws.String(BUCKET_NAME),
// Key:    aws.String(filename),
// Body:   *file,
// })
//
// return err
// }
