package entity

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type ProductCategory struct {
	Base
	ID         string `json:"id" dynamodbav:"id"`
	CategoryID string `json:"category_id" dynamodbav:"category_id"`
	Name       string `json:"name" dynamodbav:"name"`
	Image      string `json:"image" dynamodbav:"image"`
}

var _ Entity = &ProductCategory{}

func (c *ProductCategory) SetKey() Entity {
	c.PK = string(PrefixProductCategory)
	c.SK = fmt.Sprintf("%v#%v", c.Prefix(), c.ID)
	return c
}

func (c *ProductCategory) Key() Key {
	return Key{
		"pk": &types.AttributeValueMemberS{Value: c.PK},
		"sk": &types.AttributeValueMemberS{Value: c.SK},
	}
}

func (c *ProductCategory) Prefix() Prefix {
	if c.CategoryID == "" {
		return PrefixProductCategory
	}

	return PrefixProductSubCategory
}

func (c *ProductCategory) ObjectKey() string {
	if c.CategoryID == "" {
		return fmt.Sprintf("product_categories/%v", c.ID)
	}

	return fmt.Sprintf("product_categories/%v/categories/%v", c.CategoryID, c.ID)
}
