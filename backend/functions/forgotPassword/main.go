package main

import (
	"context"

	"github.com/Ekenzy-101/Serverless-Ecommerce-API/entity"
	"github.com/Ekenzy-101/Serverless-Ecommerce-API/presenter"
	"github.com/Ekenzy-101/Serverless-Ecommerce-API/usecase"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type RequestBody struct {
	Username string `json:"username" validate:"required"`
}

func Handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	svc := usecase.Default
	body := &RequestBody{}
	if err := svc.ParseRequestBody(request.Body, body); err != nil {
		return presenter.ResponseBadRequest(err.Error())
	}

	if err := svc.ValidateRequestBody(body); err != nil {
		return presenter.ResponseUnprocessableEntity(err.Error())
	}

	if err := svc.ForgotPassword(ctx, entity.AuthInput{Username: body.Username}); err != nil {
		return presenter.ResponseWrappedErrorWithCode(err)
	}

	return presenter.ResponseOK(entity.M{"message": "A confirmation code has been sent to your email"})
}

func main() {
	lambda.Start(Handler)
}
