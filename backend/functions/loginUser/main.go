package main

import (
	"context"
	"net/http"

	"github.com/Ekenzy-101/Serverless-Ecommerce-API/entity"
	"github.com/Ekenzy-101/Serverless-Ecommerce-API/presenter"
	"github.com/Ekenzy-101/Serverless-Ecommerce-API/usecase"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type RequestBody struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=8,max=99,password"`
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

	result, err := svc.LoginUser(ctx, entity.AuthInput{
		Username: body.Username,
		Password: body.Password,
	})
	if err != nil {
		return presenter.ResponseWrappedErrorWithCode(err)
	}

	if _, ok := result["session"]; ok {
		return presenter.ResponseCodeAndBody(http.StatusUnauthorized, entity.M{
			"message": "MFA verification required",
			"result":  result,
		})
	}

	return presenter.ResponseOK(result)
}

func main() {
	lambda.Start(Handler)
}
