package main

import (
	"context"

	"github.com/Ekenzy-101/Serverless-Ecommerce-API/config"
	"github.com/Ekenzy-101/Serverless-Ecommerce-API/usecase"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, event events.CognitoEventUserPoolsPreSignup) (events.CognitoEventUserPoolsPreSignup, error) {
	if config.IsTesting() {
		event.Response.AutoVerifyEmail = true
		event.Response.AutoConfirmUser = true
	}

	svc := usecase.Default
	return event, svc.IsEmailInUse(ctx, event.Request.UserAttributes["email"])
}

func main() {
	lambda.Start(Handler)
}
