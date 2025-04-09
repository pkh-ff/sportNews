package helper

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"go.uber.org/zap"
	"sportNews/pkg/log"
)

// UploadToS3
// 上傳圖片到AWS s3
func UploadToS3(s3Client *s3.Client, imgData []byte, bucketName, objectKey, contentType string) error {
	_, err := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      &bucketName,
		Key:         &objectKey,
		Body:        bytes.NewReader(imgData),
		ContentType: aws.String(contentType), // 根據圖片類型調整
		ACL:         types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		log.Error("UploadToS3: Upload file to s3 fail", zap.Error(err))
		return err
	}

	return nil
}
