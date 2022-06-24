package usecase

import (
	"context"

	"github.com/Ekenzy-101/Serverless-Ecommerce-API/entity"
)

type ProductCategoryUsecase interface {
	CreateProductCategory(ctx context.Context, category *entity.ProductCategory) (uploadURL string, err error)
	GetProductCategory(ctx context.Context, id string, categoryId string) (*entity.ProductCategory, error)
	GetProductCategories(ctx context.Context, categoryId string) ([]*entity.ProductCategory, error)
}

func (u *usecase) CreateProductCategory(ctx context.Context, category *entity.ProductCategory) (uploadURL string, err error) {
	return u.repo.CreateProductCategory(ctx, category)
}

func (u *usecase) GetProductCategory(ctx context.Context, id string, categoryId string) (*entity.ProductCategory, error) {
	return u.repo.GetProductCategory(ctx, id, categoryId)
}

func (u *usecase) GetProductCategories(ctx context.Context, categoryId string) ([]*entity.ProductCategory, error) {
	return u.repo.GetProductCategories(ctx, categoryId)
}
