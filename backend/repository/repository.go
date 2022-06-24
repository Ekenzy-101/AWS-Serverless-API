package repository

import (
	"context"
	cryptorand "crypto/rand"
	mathrand "math/rand"
	"strings"
	"time"

	"github.com/Ekenzy-101/Serverless-Ecommerce-API/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/oklog/ulid/v2"
)

type Repository interface {
	AuthRepository
	ProductCategoryRepository
	UserRepository
}

type repository struct {
	authClient *cognitoidentityprovider.Client
	dbClient   *dynamodb.Client
	objClient  *s3.PresignClient
}

func New(authClient *cognitoidentityprovider.Client, dbClient *dynamodb.Client, objClient *s3.PresignClient) Repository {
	return &repository{authClient, dbClient, objClient}
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

func (r *repository) generateUploadUrl(ctx context.Context, key string) (string, error) {
	req, err := r.objClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(config.BucketName()),
		Key:    aws.String(key),
	})
	if err != nil {
		return "", err
	}

	return req.URL, nil
}
