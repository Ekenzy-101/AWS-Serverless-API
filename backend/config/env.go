package config

import (
	"os"
	"strings"
)

func BucketName() string {
	return os.Getenv("BUCKET_NAME")
}

func CognitoAppClientID() string {
	return os.Getenv("COGNITO_APP_CLIENT_ID")
}

func CognitoAppClientSecret() string {
	return os.Getenv("COGNITO_APP_CLIENT_SECRET")
}

func CognitoUserPoolID() string {
	return os.Getenv("COGNITO_USER_POOL_ID")
}

func TableName() string {
	return os.Getenv("TABLE_NAME")
}

func IsDevelopment() bool {
	return strings.Contains(os.Getenv("APP_ENV"), "dev")
}

func IsTesting() bool {
	return strings.Contains(os.Getenv("APP_ENV"), "test")
}
