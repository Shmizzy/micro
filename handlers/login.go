package handlers

import (
	"encoding/json"
	"net/http"
	"service/models"
	"service/services"

	"github.com/aws/aws-lambda-go/events"
)

func LoginHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var credentials models.Credentials
	err := json.Unmarshal([]byte(request.Body), &credentials)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest, Body: "Invalid request body"}, nil
	}

	token, err := services.LoginUser(credentials)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusUnauthorized, Body: err.Error()}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: token}, nil
}
