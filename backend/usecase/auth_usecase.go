package usecase

import (
	"context"

	"github.com/Ekenzy-101/Serverless-Ecommerce-API/entity"
)

type AuthUsecase interface {
	ChangePassword(ctx context.Context, params entity.AuthInput) error
	DisableUserMFA(ctx context.Context, params entity.AuthInput) error
	EnableUserMFA(ctx context.Context, params entity.AuthInput) error
	ForgotPassword(ctx context.Context, params entity.AuthInput) error
	LoginUser(ctx context.Context, params entity.AuthInput) (entity.M, error)
	LogoutUser(ctx context.Context, params entity.AuthInput) error
	RefreshTokens(ctx context.Context, params entity.AuthInput) (entity.M, error)
	RegisterUser(ctx context.Context, user *entity.User) error
	ResetPassword(ctx context.Context, params entity.AuthInput) error
	SetupUserMFA(ctx context.Context, params entity.AuthInput) (string, error)
	VerifyLogin(ctx context.Context, params entity.AuthInput) (entity.M, error)
	VerifyUser(ctx context.Context, params entity.AuthInput) error
	VerifyUserMFA(ctx context.Context, params entity.AuthInput) (string, error)
}

func (u *usecase) ChangePassword(ctx context.Context, params entity.AuthInput) error {
	return u.repo.ChangePassword(ctx, params)
}

func (u *usecase) DisableUserMFA(ctx context.Context, params entity.AuthInput) error {
	params.EnableMFA = false
	return u.repo.SetUserMFAPreference(ctx, params)
}

func (u *usecase) EnableUserMFA(ctx context.Context, params entity.AuthInput) error {
	params.EnableMFA = true
	return u.repo.SetUserMFAPreference(ctx, params)
}

func (u *usecase) SetupUserMFA(ctx context.Context, params entity.AuthInput) (string, error) {
	return u.repo.SetupUserMFA(ctx, params)
}

func (u *usecase) VerifyUserMFA(ctx context.Context, params entity.AuthInput) (string, error) {
	return u.repo.VerifyUserMFA(ctx, params)
}

func (u *usecase) ForgotPassword(ctx context.Context, params entity.AuthInput) error {
	return u.repo.ForgotPassword(ctx, params)
}

func (u *usecase) LoginUser(ctx context.Context, params entity.AuthInput) (entity.M, error) {
	return u.repo.LoginUser(ctx, params)
}

func (u *usecase) LogoutUser(ctx context.Context, params entity.AuthInput) error {
	return u.repo.LogoutUser(ctx, params)
}

func (u *usecase) RegisterUser(ctx context.Context, user *entity.User) error {
	return u.repo.RegisterUser(ctx, user)
}

func (u *usecase) RefreshTokens(ctx context.Context, params entity.AuthInput) (entity.M, error) {
	return u.repo.RefreshTokens(ctx, params)
}

func (u *usecase) ResetPassword(ctx context.Context, params entity.AuthInput) error {
	return u.repo.ResetPassword(ctx, params)
}

func (u *usecase) VerifyLogin(ctx context.Context, params entity.AuthInput) (entity.M, error) {
	return u.repo.VerifyLogin(ctx, params)
}

func (u *usecase) VerifyUser(ctx context.Context, params entity.AuthInput) error {
	return u.repo.VerifyUser(ctx, params)
}
