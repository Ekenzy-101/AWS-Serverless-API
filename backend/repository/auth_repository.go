package repository

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Ekenzy-101/Serverless-Ecommerce-API/config"
	"github.com/Ekenzy-101/Serverless-Ecommerce-API/entity"
	"github.com/Ekenzy-101/Serverless-Ecommerce-API/presenter"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/aws/smithy-go"
)

type AuthRepository interface {
	ChangePassword(ctx context.Context, params entity.AuthInput) error
	ForgotPassword(ctx context.Context, params entity.AuthInput) error
	LoginUser(ctx context.Context, params entity.AuthInput) (entity.M, error)
	LogoutUser(ctx context.Context, params entity.AuthInput) error
	RefreshTokens(ctx context.Context, params entity.AuthInput) (entity.M, error)
	RegisterUser(ctx context.Context, user *entity.User) error
	ResetPassword(ctx context.Context, params entity.AuthInput) error
	SetupUserMFA(ctx context.Context, params entity.AuthInput) (string, error)
	SetUserMFAPreference(ctx context.Context, params entity.AuthInput) error
	VerifyLogin(ctx context.Context, params entity.AuthInput) (entity.M, error)
	VerifyUser(ctx context.Context, params entity.AuthInput) error
	VerifyUserMFA(ctx context.Context, params entity.AuthInput) (string, error)
}

func (r *repository) ChangePassword(ctx context.Context, params entity.AuthInput) error {
	result, err := r.authClient.ChangePassword(ctx, &cognitoidentityprovider.ChangePasswordInput{
		AccessToken:      aws.String(params.AccessToken),
		PreviousPassword: aws.String(params.Password),
		ProposedPassword: aws.String(params.NewPassword),
	})
	var ae smithy.APIError
	if errors.As(err, &ae) && ae.ErrorFault() == smithy.FaultClient {
		return presenter.WrapAPIError(ae, http.StatusUnprocessableEntity)
	}

	if err != nil {
		return presenter.WrapError(err, http.StatusInternalServerError)
	}

	r.logResult("ChangePassword", result)
	return nil
}

func (r *repository) ForgotPassword(ctx context.Context, params entity.AuthInput) error {
	result, err := r.authClient.ForgotPassword(ctx, &cognitoidentityprovider.ForgotPasswordInput{
		ClientId:   aws.String(config.CognitoAppClientID()),
		Username:   aws.String(params.Username),
		SecretHash: aws.String(r.hashUsername(params.Username)),
	})
	var ae smithy.APIError
	if errors.As(err, &ae) && ae.ErrorFault() == smithy.FaultClient {
		return presenter.WrapAPIError(ae, http.StatusUnprocessableEntity)
	}

	if err != nil {
		return presenter.WrapError(err, http.StatusInternalServerError)
	}

	r.logResult("ForgotPassword", result)
	return nil
}

func (r *repository) LoginUser(ctx context.Context, params entity.AuthInput) (entity.M, error) {
	authParams := map[string]string{
		"USERNAME":    params.Username,
		"PASSWORD":    params.Password,
		"SECRET_HASH": r.hashUsername(params.Username),
	}
	result, err := r.initiateAuth(ctx, types.AuthFlowTypeUserPasswordAuth, authParams)
	if err != nil {
		return nil, err
	}

	r.logResult("LoginUser", result)
	if result.AuthenticationResult == nil {
		return entity.M{
			"session":              result.Session,
			"challenge_parameters": result.ChallengeParameters,
		}, nil
	}

	return r.modelAuthResult(result.AuthenticationResult), nil
}

func (r *repository) LogoutUser(ctx context.Context, params entity.AuthInput) error {
	result, err := r.authClient.RevokeToken(ctx, &cognitoidentityprovider.RevokeTokenInput{
		ClientId:     aws.String(config.CognitoAppClientID()),
		ClientSecret: aws.String(config.CognitoAppClientSecret()),
		Token:        aws.String(params.RefreshToken),
	})
	var ae smithy.APIError
	if errors.As(err, &ae) && ae.ErrorFault() == smithy.FaultClient {
		return presenter.WrapAPIError(ae, http.StatusUnprocessableEntity)
	}

	if err != nil {
		return presenter.WrapError(err, http.StatusInternalServerError)
	}

	r.logResult("LogoutUser", result)
	return nil
}

func (r *repository) RefreshTokens(ctx context.Context, params entity.AuthInput) (entity.M, error) {
	authParams := map[string]string{
		"REFRESH_TOKEN": params.RefreshToken,
		"SECRET_HASH":   r.hashUsername(params.Username),
	}
	result, err := r.initiateAuth(ctx, types.AuthFlowTypeRefreshTokenAuth, authParams)
	if err != nil {
		return nil, err
	}

	r.logResult("RefreshTokens", result)
	return r.modelAuthResult(result.AuthenticationResult), nil
}

func (r *repository) RegisterUser(ctx context.Context, user *entity.User) (err error) {
	user.ID, err = r.generateId()
	if err != nil {
		return presenter.WrapError(err, http.StatusInternalServerError)
	}

	result, err := r.authClient.SignUp(ctx, &cognitoidentityprovider.SignUpInput{
		Password:   aws.String(user.Password),
		Username:   aws.String(user.ID),
		ClientId:   aws.String(config.CognitoAppClientID()),
		SecretHash: aws.String(r.hashUsername(user.ID)),
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("name"),
				Value: aws.String(user.Name),
			},
			{
				Name:  aws.String("email"),
				Value: aws.String(user.Email),
			},
		},
	})
	var ae smithy.APIError
	if errors.As(err, &ae) && ae.ErrorFault() == smithy.FaultClient {
		return presenter.WrapAPIError(ae, http.StatusUnprocessableEntity)
	}

	if err != nil {
		return presenter.WrapError(err, http.StatusInternalServerError)
	}

	r.logResult("RegisterUser", result)
	return nil
}

func (r *repository) ResetPassword(ctx context.Context, params entity.AuthInput) error {
	result, err := r.authClient.ConfirmForgotPassword(ctx, &cognitoidentityprovider.ConfirmForgotPasswordInput{
		Password:         aws.String(params.Password),
		ConfirmationCode: aws.String(params.Code),
		ClientId:         aws.String(config.CognitoAppClientID()),
		Username:         aws.String(params.Username),
		SecretHash:       aws.String(r.hashUsername(params.Username)),
	})
	var ae smithy.APIError
	if errors.As(err, &ae) && ae.ErrorFault() == smithy.FaultClient {
		return presenter.WrapAPIError(ae, http.StatusUnprocessableEntity)
	}

	if err != nil {
		return presenter.WrapError(err, http.StatusInternalServerError)
	}

	r.logResult("ResetPassword", result)
	return nil
}

func (r *repository) SetupUserMFA(ctx context.Context, params entity.AuthInput) (string, error) {
	result, err := r.authClient.AssociateSoftwareToken(ctx, &cognitoidentityprovider.AssociateSoftwareTokenInput{
		AccessToken: aws.String(params.AccessToken),
	})
	var ae smithy.APIError
	if errors.As(err, &ae) && ae.ErrorFault() == smithy.FaultClient {
		return "", presenter.WrapAPIError(ae, http.StatusUnprocessableEntity)
	}

	if err != nil {
		return "", presenter.WrapError(err, http.StatusInternalServerError)
	}

	r.logResult("SetUserMFA", result)
	return *result.SecretCode, nil
}

func (r *repository) SetUserMFAPreference(ctx context.Context, params entity.AuthInput) error {
	result, err := r.authClient.SetUserMFAPreference(ctx, &cognitoidentityprovider.SetUserMFAPreferenceInput{
		AccessToken: aws.String(params.AccessToken),
		SoftwareTokenMfaSettings: &types.SoftwareTokenMfaSettingsType{
			Enabled: params.EnableMFA,
		},
	})
	var ae smithy.APIError
	if errors.As(err, &ae) && ae.ErrorFault() == smithy.FaultClient {
		return presenter.WrapAPIError(ae, http.StatusUnprocessableEntity)
	}

	if err != nil {
		return presenter.WrapError(err, http.StatusInternalServerError)
	}

	r.logResult("SetUserMFA", result)
	return nil
}

func (r *repository) VerifyLogin(ctx context.Context, params entity.AuthInput) (entity.M, error) {
	result, err := r.authClient.RespondToAuthChallenge(ctx, &cognitoidentityprovider.RespondToAuthChallengeInput{
		ChallengeName: types.ChallengeNameTypeSoftwareTokenMfa,
		ClientId:      aws.String(config.CognitoAppClientID()),
		Session:       aws.String(params.Session),
		ChallengeResponses: map[string]string{
			"USERNAME":                params.Username,
			"SOFTWARE_TOKEN_MFA_CODE": params.Code,
			"SECRET_HASH":             r.hashUsername(params.Username),
		},
	})
	var ae smithy.APIError
	if errors.As(err, &ae) && ae.ErrorFault() == smithy.FaultClient {
		return nil, presenter.WrapAPIError(ae, http.StatusUnprocessableEntity)
	}

	if err != nil {
		return nil, presenter.WrapError(err, http.StatusInternalServerError)
	}

	r.logResult("VerifyLoginUser", result)
	return r.modelAuthResult(result.AuthenticationResult), nil
}

func (r *repository) VerifyUser(ctx context.Context, params entity.AuthInput) error {
	result, err := r.authClient.ConfirmSignUp(ctx, &cognitoidentityprovider.ConfirmSignUpInput{
		ConfirmationCode: aws.String(params.Code),
		Username:         aws.String(params.Username),
		ClientId:         aws.String(config.CognitoAppClientID()),
		SecretHash:       aws.String(r.hashUsername(params.Username)),
	})
	var ae smithy.APIError
	if errors.As(err, &ae) && ae.ErrorFault() == smithy.FaultClient {
		return presenter.WrapAPIError(ae, http.StatusUnprocessableEntity)
	}

	if err != nil {
		return presenter.WrapError(err, http.StatusInternalServerError)
	}

	r.logResult("VerifyUser", result)
	return nil
}

func (r *repository) VerifyUserMFA(ctx context.Context, params entity.AuthInput) (string, error) {
	result, err := r.authClient.VerifySoftwareToken(ctx, &cognitoidentityprovider.VerifySoftwareTokenInput{
		UserCode:    aws.String(params.Code),
		AccessToken: aws.String(params.AccessToken),
	})
	var ae smithy.APIError
	if errors.As(err, &ae) && ae.ErrorFault() == smithy.FaultClient {
		return "", presenter.WrapAPIError(ae, http.StatusUnprocessableEntity)
	}

	if err != nil {
		return "", presenter.WrapError(err, http.StatusInternalServerError)
	}

	r.logResult("VerifyUserMFA", result)
	return string(result.Status), nil
}

func (r *repository) hashUsername(username string) string {
	mac := hmac.New(sha256.New, []byte(config.CognitoAppClientSecret()))
	mac.Write([]byte(username + config.CognitoAppClientID()))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (r *repository) logResult(prefix string, result interface{}) {
	payload, _ := json.MarshalIndent(result, "", "  ")
	log.Println(prefix, string(payload))
}

func (r *repository) initiateAuth(ctx context.Context, authFlow types.AuthFlowType, authParams map[string]string) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	result, err := r.authClient.InitiateAuth(ctx, &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow:       authFlow,
		ClientId:       aws.String(config.CognitoAppClientID()),
		AuthParameters: authParams,
	})

	var ae smithy.APIError
	if errors.As(err, &ae) && ae.ErrorFault() == smithy.FaultClient {
		return nil, presenter.WrapAPIError(ae, http.StatusUnprocessableEntity)
	}

	if err != nil {
		return nil, presenter.WrapError(err, http.StatusInternalServerError)
	}

	return result, nil
}

func (r *repository) modelAuthResult(result *types.AuthenticationResultType) entity.M {
	return entity.M{
		"access_token":  result.AccessToken,
		"id_token":      result.IdToken,
		"refresh_token": result.RefreshToken,
		"token_type":    result.TokenType,
	}
}
