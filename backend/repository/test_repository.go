package repository

import (
	"context"

	"github.com/Ekenzy-101/Serverless-Ecommerce-API/config"
	"github.com/Ekenzy-101/Serverless-Ecommerce-API/entity"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type TestRepository interface {
	DeleteUser(ctx context.Context, id string) error
}

func (r *repository) DeleteUser(ctx context.Context, id string) error {
	_, err := r.authClient.AdminDeleteUser(ctx, &cognitoidentityprovider.AdminDeleteUserInput{
		Username:   aws.String(id),
		UserPoolId: aws.String(config.CognitoUserPoolID()),
	})
	if err != nil {
		return err
	}

	_, err = r.dbClient.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key:       (&entity.User{ID: id}).SetKey().Key(),
		TableName: aws.String(config.TableName()),
	})
	return err
}
