package presenter

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func ResponseBadRequest(msg string) (events.APIGatewayProxyResponse, error) {
	return ResponseCodeAndBody(http.StatusBadRequest, map[string]string{"message": msg})
}

func ResponseCodeAndBody(code int, body interface{}) (events.APIGatewayProxyResponse, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Body:       string(payload),
		Headers: map[string]string{
			"Content-type": "application/json",
		},
	}, nil
}

func ResponseCreated(body interface{}) (events.APIGatewayProxyResponse, error) {
	return ResponseCodeAndBody(http.StatusCreated, body)
}

func ResponseInternalServerError() (events.APIGatewayProxyResponse, error) {
	return ResponseCodeAndBody(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
}

func ResponseNoContent() (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusNoContent,
		Headers: map[string]string{
			"Content-type": "application/json",
		},
	}, nil
}

func ResponseNotFound(msg string) (events.APIGatewayProxyResponse, error) {
	return ResponseCodeAndBody(http.StatusNotFound, map[string]string{"message": msg})
}

func ResponseOK(body interface{}) (events.APIGatewayProxyResponse, error) {
	return ResponseCodeAndBody(http.StatusOK, body)
}

func ResponseUnprocessableEntity(msg string) (events.APIGatewayProxyResponse, error) {
	return ResponseCodeAndBody(http.StatusUnprocessableEntity, map[string]string{"message": msg})
}

func ResponseWrappedErrorWithCode(err error) (events.APIGatewayProxyResponse, error) {
	code, err := UnwrapCodeError(err)
	if code == http.StatusInternalServerError {
		log.Println("Unexpected error", err)
		return ResponseInternalServerError()
	}

	return ResponseCodeAndBody(code, map[string]string{"message": err.Error()})
}
