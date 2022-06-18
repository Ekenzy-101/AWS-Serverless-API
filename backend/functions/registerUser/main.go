package main

import (
	"context"
	"log"

	"github.com/Ekenzy-101/Serverless-Ecommerce-API/entity"
	"github.com/Ekenzy-101/Serverless-Ecommerce-API/presenter"
	"github.com/Ekenzy-101/Serverless-Ecommerce-API/usecase"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type RequestBody struct {
	Email    string `json:"email" validate:"required,max=255,email"`
	Password string `json:"password" validate:"required,min=8,max=99,password"`
	Name     string `json:"name" validate:"required,max=100,name"`
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

	user := &entity.User{
		Email:    body.Email,
		Name:     body.Name,
		Password: body.Password,
	}
	if err := svc.RegisterUser(ctx, user); err != nil {
		return presenter.ResponseWrappedErrorWithCode(err)
	}

	if err := svc.CreateUser(ctx, user); err != nil {
		log.Println(err)
		return presenter.ResponseInternalServerError()
	}

	return presenter.ResponseCreated(entity.M{"user": user})
}

func main() {
	lambda.Start(Handler)
}
