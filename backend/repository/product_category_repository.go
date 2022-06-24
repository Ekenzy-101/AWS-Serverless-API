package repository

import (
	"context"
	"strings"

	"github.com/Ekenzy-101/Serverless-Ecommerce-API/config"
	"github.com/Ekenzy-101/Serverless-Ecommerce-API/entity"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type ProductCategoryRepository interface {
	GetProductCategories(ctx context.Context, categoryId string) ([]*entity.ProductCategory, error)
	CreateProductCategory(ctx context.Context, category *entity.ProductCategory) (uploadURL string, err error)
	GetProductCategory(ctx context.Context, id string, categoryId string) (*entity.ProductCategory, error)
}

func (r *repository) CreateProductCategory(ctx context.Context, category *entity.ProductCategory) (uploadURL string, err error) {
	category.ID, err = r.generateId()
	if err != nil {
		return
	}

	uploadURL, err = r.generateUploadUrl(ctx, category.ObjectKey())
	if err != nil {
		return
	}

	category.Image = strings.Split(uploadURL, "?")[0]
	item, err := attributevalue.MarshalMap(category.SetKey())
	if err != nil {
		return "", err
	}

	_, err = r.dbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(config.TableName()),
		Item:      item,
	})
	return
}

func (r *repository) GetProductCategory(ctx context.Context, id string, categoryId string) (*entity.ProductCategory, error) {
	category := &entity.ProductCategory{ID: id, CategoryID: categoryId}
	result, err := r.dbClient.GetItem(ctx, &dynamodb.GetItemInput{
		Key:       category.SetKey().Key(),
		TableName: aws.String(config.TableName()),
	})
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, nil
	}

	return category, attributevalue.UnmarshalMap(result.Item, category)
}

func (r *repository) GetProductCategories(ctx context.Context, categoryId string) ([]*entity.ProductCategory, error) {
	var filterExpression *string
	expressionAttributeNames := map[string]string{"#pk": "pk"}
	expressionAttributeValues := map[string]types.AttributeValue{
		":pk": &types.AttributeValueMemberS{Value: string(entity.PrefixProductCategory)},
	}
	keyConditionExpression := "#pk = :pk"
	if categoryId != "" {
		expressionAttributeNames["#category_id"] = "category_id"
		expressionAttributeNames["#id"] = "id"
		expressionAttributeValues[":category_id"] = &types.AttributeValueMemberS{Value: categoryId}
		filterExpression = aws.String("#id = :category_id or #category_id = :category_id")
	}

	result, err := r.dbClient.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(config.TableName()),
		KeyConditionExpression:    aws.String(keyConditionExpression),
		ExpressionAttributeNames:  expressionAttributeNames,
		ExpressionAttributeValues: expressionAttributeValues,
		FilterExpression:          filterExpression,
	})
	if err != nil {
		return nil, err
	}

	categories := make([]*entity.ProductCategory, 0, result.Count)
	return categories, attributevalue.UnmarshalListOfMaps(result.Items, &categories)
}
