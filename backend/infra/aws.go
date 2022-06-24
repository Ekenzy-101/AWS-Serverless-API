package infra

import (
	"context"
	"log"

	appconfig "github.com/Ekenzy-101/Serverless-Ecommerce-API/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var cfg aws.Config

func init() {
	var err error

	opts := []func(*config.LoadOptions) error{}
	if appconfig.IsDevelopment() {
		opts = append(opts, config.WithClientLogMode(aws.LogRetries|aws.LogRequest))
	}

	cfg, err = config.LoadDefaultConfig(context.TODO(), opts...)
	if err != nil {
		log.Fatal(err)
	}
}

func NewDatabaseClient() *dynamodb.Client {
	return dynamodb.NewFromConfig(cfg)
}

func NewAuthClient() *cognitoidentityprovider.Client {
	return cognitoidentityprovider.NewFromConfig(cfg)
}

func NewObjectClient() *s3.PresignClient {
	return s3.NewPresignClient(s3.NewFromConfig(cfg))
}
