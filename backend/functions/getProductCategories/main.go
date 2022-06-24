package main

import (
	"context"
	"log"

	"github.com/Ekenzy-101/Serverless-Ecommerce-API/presenter"
	"github.com/Ekenzy-101/Serverless-Ecommerce-API/usecase"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	svc := usecase.Default
	categories, err := svc.GetProductCategories(ctx, request.QueryStringParameters["category_id"])
	if err != nil {
		log.Println("get product categories", err)
		return presenter.ResponseInternalServerError()
	}

	return presenter.ResponseOK(categories)
}

func main() {
	lambda.Start(Handler)
}
