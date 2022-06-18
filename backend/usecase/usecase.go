package usecase

import (
	"github.com/Ekenzy-101/Serverless-Ecommerce-API/infra"
	"github.com/Ekenzy-101/Serverless-Ecommerce-API/repository"
)

type Usecase interface {
	AuthUsecase
	UserUsecase
	RequestUsecase
}

type usecase struct {
	repo repository.Repository
}

func init() {
	repo := repository.New(infra.NewAuthClient(), infra.NewDatabaseClient())
	Default = New(repo)
}

var Default Usecase

func New(repo repository.Repository) Usecase {
	return &usecase{repo}
}
