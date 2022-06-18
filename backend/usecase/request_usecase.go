package usecase

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-playground/validator/v10"
)

type RequestUsecase interface {
	ExtractTokenFromRequest(req interface{}) string
	ExtractUsernameFromRequest(req interface{}) string
	ParseRequestBody(body string, obj interface{}) error
	ValidateRequestBody(obj interface{}) error
}

var (
	lowercaseRegex        = regexp.MustCompile(`[a-z]+`)
	nameRegex             = regexp.MustCompile(`^[a-zA-z ]+$`)
	numberRegex           = regexp.MustCompile(`\d+`)
	specialCharacterRegex = regexp.MustCompile(`\W+`)
	uppercaseRegex        = regexp.MustCompile(`[A-Z]+`)

	validate *validator.Validate
)

func init() {
	validate = validator.New()
	if err := validate.RegisterValidation("name", validateName); err != nil {
		log.Panicf("RegisterValidation [name] %v", err)
	}

	if err := validate.RegisterValidation("password", validatePassword); err != nil {
		log.Fatalf("RegisterValidation [password] %v", err)
	}
}

func (u *usecase) ExtractTokenFromRequest(req interface{}) string {
	request, ok := req.(events.APIGatewayV2HTTPRequest)
	if !ok {
		return ""
	}

	token := request.Headers["authorization"]
	return strings.Replace(token, "Bearer ", "", 1)
}

func (u *usecase) ExtractUsernameFromRequest(req interface{}) string {
	request, ok := req.(events.APIGatewayV2HTTPRequest)
	if !ok {
		return ""
	}

	cliams := request.RequestContext.Authorizer.JWT.Claims
	if cliams["token_use"] == "access" {
		return cliams["username"]
	}

	return cliams["cognito:username"]
}

func (u *usecase) ParseRequestBody(body string, obj interface{}) error {
	return json.Unmarshal([]byte(body), obj)
}

func (u *usecase) ValidateRequestBody(obj interface{}) error {
	err := validate.Struct(obj)
	fieldErrors, ok := err.(validator.ValidationErrors)
	if err != nil && !ok {
		return err
	}

	if err != nil {
		return transformFieldError(fieldErrors[0])
	}

	return nil
}

func transformFieldError(err validator.FieldError) error {
	switch err.ActualTag() {
	case "name":
		return fmt.Errorf("%v should contain only letters and spaces", err.StructField())
	case "email":
		return fmt.Errorf("%v is not a valid email address", err.StructField())
	case "gt":
		return fmt.Errorf("%v should be greater than %v", err.StructField(), err.Param())
	case "lte":
		return fmt.Errorf("%v should not be greater than %v", err.StructField(), err.Param())
	case "max":
		return fmt.Errorf("%v should be at most %v characters", err.StructField(), err.Param())
	case "min":
		return fmt.Errorf("%v should be at least %v characters", err.StructField(), err.Param())
	case "number":
		return fmt.Errorf("%v should contain only numbers", err.StructField())
	case "password":
		return fmt.Errorf("%v should contain at least 1 uppercase, lowercase, numeric and special characters", err.StructField())
	case "required":
		return fmt.Errorf("%v is required", err.StructField())
	default:
		return fmt.Errorf("%v is in an invalid format", err.StructField())
	}
}

func validateName(fl validator.FieldLevel) bool {
	return nameRegex.MatchString(fl.Field().String())
}

func validatePassword(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return specialCharacterRegex.MatchString(value) &&
		lowercaseRegex.MatchString(value) &&
		uppercaseRegex.MatchString(value) &&
		numberRegex.MatchString(value)
}
