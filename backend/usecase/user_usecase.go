package usecase

import (
	"context"

	"github.com/Ekenzy-101/Serverless-Ecommerce-API/entity"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, user *entity.User) error
}

func (u *usecase) CreateUser(ctx context.Context, user *entity.User) error {
	return u.repo.CreateUser(ctx, user)
}
