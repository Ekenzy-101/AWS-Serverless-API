package entity

import "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

type Base struct {
	PK string `json:"-" dynamodbav:"pk"`
	SK string `json:"-" dynamodbav:"sk"`
}

type Key map[string]types.AttributeValue

// Shortcut for map[string]interface{}
type M map[string]interface{}

type AuthInput struct {
	AccessToken  string
	Code         string
	Password     string
	NewPassword  string
	RefreshToken string
	Session      string
	Username     string
	EnableMFA    bool
}

type Entity interface {
	SetKey() Entity
	Key() Key
}
