package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

func main() {
	args := os.Args[1:]

	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	img := path + "/" + args[0]
	s3URL, err := uploadImgToS3(img)
	if err != nil {
		panic(err)
	}
	fmt.Println(s3URL)

}

const REGION = "ap-northeast-1"
const BUCKET = "ad-prodjp-testing-data"

func uploadImgToS3(path string) (string, error) {
	// Open the file for use
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Get file size and read the file content into a buffer
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	contentType := http.DetectContentType(buffer)

	key := uuid.New().String() + "_" + filepath.Base(path)

	session, _ := session.NewSession(aws.NewConfig().WithRegion(REGION))
	s3Session := s3.New(session)
	_, err = s3Session.PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(BUCKET),
		Key:                  aws.String(key),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(contentType),
		ContentDisposition:   aws.String("inline"),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		return "", err
	}

	req, _ := s3Session.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(BUCKET),
		Key:    aws.String(key),
	})

	url, _ := req.Presign(24 * time.Hour)
	return url, nil
}
