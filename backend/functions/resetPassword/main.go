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
	Password string `json:"password" validate:"required,min=8,max=99,password"`
	Code     string `json:"code" validate:"required,number"`
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
		Username: body.Username,
		Password: body.Password,
		Code:     body.Code,
	}
	if err := svc.ResetPassword(ctx, params); err != nil {
		return presenter.ResponseWrappedErrorWithCode(err)
	}

	return presenter.ResponseOK(entity.M{"message": "Password has been resetted successfully"})
}

func main() {
	lambda.Start(Handler)
}
