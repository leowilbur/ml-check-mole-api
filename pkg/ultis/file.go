package ultis

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	bucketName  string = "molepatrol-photos"
	uploadLevel string = "private"
	awsRegion   string = "ap-southeast-2"
)

// AWSFileContent ...
type AWSFileContent struct {
	AccountID  string
	FileName   string
	FileEncode string
}

// UploadFileToAWS ...
func (content *AWSFileContent) UploadFileToAWS() error {

	regrexGroup := regexp.MustCompile(`^data:(image\/\w+);base64,(.+)$`).FindStringSubmatch(content.FileEncode)
	if len(regrexGroup) != 3 {
		return errors.New("Invalid image data base64")
	}
	supportFilContentList := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
	}

	if !supportFilContentList[regrexGroup[1]] {
		return errors.New("Not support content-type: " + regrexGroup[1])
	}

	dec, err := base64.StdEncoding.DecodeString(regrexGroup[2])
	if err != nil {
		return errors.New("Decode file content fail")
	}
	size := int64(len(dec))

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion)},
	)

	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(bucketName),
		Key:                  aws.String(fmt.Sprintf("%s/%s:%s/%s", uploadLevel, awsRegion, content.AccountID, content.FileName)),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(dec),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(regrexGroup[1]),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		return errors.New("Can not upload file to AWS S3")
	}
	return nil
}

// DownloadFileFromAWS ...
func (content *AWSFileContent) DownloadFileFromAWS() (string, error) {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion)},
	)
	svc := s3.New(sess)

	params := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fmt.Sprintf("%s/%s:%s/%s", uploadLevel, awsRegion, content.AccountID, content.FileName)),
	}

	req, _ := svc.GetObjectRequest(params)

	url, _ := req.Presign(24 * time.Hour) // expiration time is 60 mins
	fmt.Println(url)
	return url, nil
}
