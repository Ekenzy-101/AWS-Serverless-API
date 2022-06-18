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
	OldPassword string `json:"old_password" validate:"required,password"`
	NewPassword string `json:"new_password" validate:"required,password"`
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
		Password:    body.OldPassword,
		NewPassword: body.NewPassword,
		AccessToken: svc.ExtractTokenFromRequest(request),
	}
	if err := svc.ChangePassword(ctx, params); err != nil {
		return presenter.ResponseWrappedErrorWithCode(err)
	}

	return presenter.ResponseOK(entity.M{"message": "Password has been changed successfully"})
}

func main() {
	lambda.Start(Handler)
}
