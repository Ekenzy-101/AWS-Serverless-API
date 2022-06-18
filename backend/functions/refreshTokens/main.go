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
	Username     string `json:"username" validate:"required"`
	RefreshToken string `json:"refresh_token" validate:"required"`
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

	result, err := svc.RefreshTokens(ctx, entity.AuthInput{
		Username:     body.Username,
		RefreshToken: body.RefreshToken,
	})
	if err != nil {
		return presenter.ResponseWrappedErrorWithCode(err)
	}

	return presenter.ResponseOK(result)
}

func main() {
	lambda.Start(Handler)
}
