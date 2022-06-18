package repository

import (
	cryptorand "crypto/rand"
	mathrand "math/rand"
	"strings"
	"time"

	"github.com/Ekenzy-101/Serverless-Ecommerce-API/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/oklog/ulid/v2"
)

type Repository interface {
	AuthRepository
	UserRepository
}

type repository struct {
	authClient *cognitoidentityprovider.Client
	dbClient   *dynamodb.Client
}

func New(authClient *cognitoidentityprovider.Client, dbClient *dynamodb.Client) Repository {
	return &repository{authClient, dbClient}
}

func (r *repository) generateId() (string, error) {
	entropy := cryptorand.Reader
	if config.IsTesting() || config.IsDevelopment() {
		seed := time.Now().UTC().UnixNano()
		source := mathrand.NewSource(seed)
		entropy = mathrand.New(source)
	}

	id, err := ulid.New(ulid.Now(), entropy)
	return strings.ToLower(id.String()), err
}
