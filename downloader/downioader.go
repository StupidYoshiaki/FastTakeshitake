package downloader

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var S3FilePaths = make(map[string]string)

var (
	s3Client    *s3.S3
	bucketName  string
	downloadDir = "./img"
)

func Init() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
	}))
	s3Client = s3.New(sess)
	bucketName = os.Getenv("AWS_S3_BUCKET")

	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}
}

func DownloadFiles() error {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	}

	err := s3Client.ListObjectsV2Pages(input, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		// cnt := 0 // test

		for _, obj := range page.Contents {
			key := *obj.Key

			log.Printf("Downloading %s", key)
			baseName := filepath.Base(key)
			localPath := filepath.Join(downloadDir, baseName)
			if err := downloadFile(key, localPath); err != nil {
				log.Printf("Failed to download file %s: %v", key, err)
				return false
			}
			// key = strings.SplitAfter(baseName, "/")[1]
			S3FilePaths[baseName] = localPath

			// cnt++ // test
			// if cnt > 2 {
			// 	break
			// }
		}
		return true
	})
	if err != nil {
		log.Printf("Failed to list objects: %v", err)
		return err
	}
	log.Printf("Downloaded %d files", len(S3FilePaths))
	return nil
}

func downloadFile(key, localPath string) error {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}

	obj, err := s3Client.GetObject(input)
	if err != nil {
		return err
	}
	defer obj.Body.Close()

	outFile, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, obj.Body)
	return err
}
