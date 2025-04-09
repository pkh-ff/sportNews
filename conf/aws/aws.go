package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"sportNews/conf"
	"sportNews/pkg/log"
)

func New(conf conf.Aws) (*s3.Client, error) {
	cfg, err := setConfig(conf)
	if err != nil {
		return nil, err
	}
	return s3.NewFromConfig(cfg), nil
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
