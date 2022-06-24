package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Ekenzy-101/Serverless-Ecommerce-API/entity"
	"github.com/Ekenzy-101/Serverless-Ecommerce-API/presenter"
	"github.com/Ekenzy-101/Serverless-Ecommerce-API/usecase"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type RequestBody struct {
	Name       string `json:"name" validate:"required"`
	CategoryID string `json:"category_id"`
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

	if body.CategoryID != "" {
		category, err := svc.GetProductCategory(ctx, body.CategoryID, "")
		if err != nil {
			log.Println("get product category", err)
			return presenter.ResponseInternalServerError()
		}

		if category == nil {
			return presenter.ResponseNotFound(fmt.Sprintf("Product category '%v' not found", body.CategoryID))
		}
	}

	category := &entity.ProductCategory{
		Name:       body.Name,
		CategoryID: body.CategoryID,
	}
	uploadURL, err := svc.CreateProductCategory(ctx, category)
	if err != nil {
		log.Println("create product category", err)
		return presenter.ResponseInternalServerError()
	}

	return presenter.ResponseCreated(entity.M{
		"category":  category,
		"uploadURL": uploadURL,
	})
}

func main() {
	lambda.Start(Handler)
}
