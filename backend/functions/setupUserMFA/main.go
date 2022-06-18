package main

import (
	"context"

	"github.com/Ekenzy-101/Serverless-Ecommerce-API/entity"
	"github.com/Ekenzy-101/Serverless-Ecommerce-API/presenter"
	"github.com/Ekenzy-101/Serverless-Ecommerce-API/usecase"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	svc := usecase.Default
	setupKey, err := svc.SetupUserMFA(ctx, entity.AuthInput{
		AccessToken: svc.ExtractTokenFromRequest(request),
	})
	if err != nil {
		return presenter.ResponseWrappedErrorWithCode(err)
	}

	return presenter.ResponseCreated(entity.M{"setup_key": setupKey})
}

func main() {
	lambda.Start(Handler)
}
