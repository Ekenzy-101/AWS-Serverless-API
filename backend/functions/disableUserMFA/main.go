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
	Code string `json:"code" validate:"required,number"`
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

	params := entity.AuthInput{
		AccessToken: svc.ExtractTokenFromRequest(request),
		Code:        body.Code,
	}
	status, err := svc.VerifyUserMFA(ctx, params)
	if err != nil {
		return presenter.ResponseWrappedErrorWithCode(err)
	}

	if err := svc.DisableUserMFA(ctx, params); err != nil {
		return presenter.ResponseWrappedErrorWithCode(err)
	}

	return presenter.ResponseOK(entity.M{"status": status})
}

func main() {
	lambda.Start(Handler)
}
