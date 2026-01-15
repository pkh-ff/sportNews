package aws

import (
	"context"
	"sportNews/conf"
	"sportNews/pkg/log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Api interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

type S3Client struct {
	Client S3Api
	Bucket string
	Acl    bool
}

func New(conf conf.Aws) (*S3Client, error) {
	cfg, err := setConfig(conf)
	if err != nil {
		return nil, err
	}

	return &S3Client{
		Client: s3.NewFromConfig(cfg),
		Bucket: conf.Bucket,
		Acl:    conf.Acl,
	}, nil
}

func setConfig(conf conf.Aws) (aws.Config, error) {
	customCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(conf.Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(conf.AccessKey, conf.SecretKey, ""),
		),
	)
	if err != nil {
		log.Errorf("failed to load config: %v", err)
		return aws.Config{}, err
	}

	return customCfg, nil
}
