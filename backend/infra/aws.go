package infra

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var cfg aws.Config

func init() {
	var err error
	cfg, err = config.LoadDefaultConfig(context.TODO(), config.WithClientLogMode(aws.LogRetries|aws.LogRequest))
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
