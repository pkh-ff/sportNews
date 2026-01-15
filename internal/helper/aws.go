package helper

import (
	"bytes"
	"context"
	confAws "sportNews/conf/aws"
	"sportNews/pkg/log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"go.uber.org/zap"
)

const (
	basePath              = "n8/"
	NewsCoverPrefix       = basePath + "news/source/"
	NewsCustomCoverPrefix = basePath + "news/"
	VideoCoverPrefix      = basePath + "video/"
	TeamIconPrefix        = basePath + "team/"
)

// UploadToS3
// 上傳圖片到AWS s3
func UploadToS3(s3Client *confAws.S3Client, imgData []byte, objectKey, contentType string) error {
	pi := s3.PutObjectInput{
		Bucket:      &s3Client.Bucket,
		Key:         &objectKey,
		Body:        bytes.NewReader(imgData),
		ContentType: aws.String(contentType), // 根據圖片類型調整
	}

	if s3Client.Acl {
		pi.ACL = types.ObjectCannedACLPublicRead
	}

	_, err := s3Client.Client.PutObject(context.TODO(), &pi)
	if err != nil {
		log.Error("UploadToS3: Upload file to s3 fail", zap.Error(err))
		return err
	}

	return nil
}
