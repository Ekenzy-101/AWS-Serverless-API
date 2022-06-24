package test

import (
	"context"
	"os"
	"testing"

	"github.com/Ekenzy-101/Serverless-Ecommerce-API/infra"
	"github.com/Ekenzy-101/Serverless-Ecommerce-API/repository"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEMAAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "EMA API")
}

var (
	baseUrl = os.Getenv("API_ENDPOINT")
	repo    repository.Repository
	ctx     context.Context
)

var _ = BeforeSuite(func() {
	By("Initialzing global variables")
	ctx = context.Background()
	repo = repository.New(infra.NewAuthClient(), infra.NewDatabaseClient(), infra.NewObjectClient())
})
