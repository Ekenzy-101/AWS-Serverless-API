package entity

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type User struct {
	Base     `json:"-" dynamodbav:",inline"`
	ID       string `json:"id,omitempty" dynamodbav:"id,omitempty"`
	Email    string `json:"email,omitempty" dynamodbav:"email,omitempty"`
	Name     string `json:"name,omitempty" dynamodbav:"name,omitempty"`
	Password string `json:"-" dynamodbav:"-"`
}

var _ Entity = &User{}

func (u *User) SetKey() Entity {
	u.PK = fmt.Sprintf("USER#%v", u.ID)
	u.SK = fmt.Sprintf("USER#%v", u.ID)
	return u
}

func (u *User) Key() Key {
	return Key{
		"pk": &types.AttributeValueMemberS{Value: u.PK},
		"sk": &types.AttributeValueMemberS{Value: u.SK},
	}
}
