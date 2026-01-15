package helper

import (
	"context"
	"errors"
	"io"
	confAws "sportNews/conf/aws"
	"sportNews/conf/aws/mocks"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUploadToS3_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockS3 := mocks.NewMockS3Api(ctrl)

	bucketName := "test-bucket"
	objectKey := "test/image.png"
	contentType := "image/png"
	imgData := []byte("fake-image-data")

	client := &confAws.S3Client{
		Client: mockS3,
		Bucket: bucketName,
		Acl:    true,
	}

	mockS3.EXPECT().
		PutObject(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, input *s3.PutObjectInput, opts ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
			// 驗證參數
			assert.Equal(t, bucketName, *input.Bucket)
			assert.Equal(t, objectKey, *input.Key)
			assert.Equal(t, contentType, *input.ContentType)
			assert.Equal(t, types.ObjectCannedACLPublicRead, input.ACL)

			// 驗證 Body 內容
			body, _ := io.ReadAll(input.Body)
			assert.Equal(t, imgData, body)

			return &s3.PutObjectOutput{}, nil
		}).
		Times(1)

	err := UploadToS3(client, imgData, objectKey, contentType)

	require.NoError(t, err)
}

func TestUploadToS3WithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockS3 := mocks.NewMockS3Api(ctrl)
	client := &confAws.S3Client{
		Client: mockS3,
		Bucket: "bucket",
		Acl:    false,
	}

	mockS3.EXPECT().
		PutObject(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("s3 error")).
		Times(1)

	err := UploadToS3(client, []byte("data"), "key", "type")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "s3 error")
}
