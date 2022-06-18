package repository

import (
	"context"

	"github.com/Ekenzy-101/Serverless-Ecommerce-API/config"
	"github.com/Ekenzy-101/Serverless-Ecommerce-API/entity"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
}

func (r *repository) CreateUser(ctx context.Context, user *entity.User) error {
	item, err := attributevalue.MarshalMap(user.SetKey())
	if err != nil {
		return err
	}

	_, err = r.dbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(config.TableName()),
		Item:      item,
	})
	return err
}
