package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"service/models"
	"service/services"

	"github.com/aws/aws-lambda-go/events"
)

func ConfirmHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Received request body: %s", request.Body)
	var confirmReq models.ConfirmUser
	err := json.Unmarshal([]byte(request.Body), &confirmReq)
	if err != nil {
		log.Printf("Error parsing request body: %v", err)
		log.Printf("Request body: %s", request.Body)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest, Body: "Invalid request body"}, nil
	}

	err = services.ConfirmUser(confirmReq.Username, confirmReq.ConfirmationCode)
	if err != nil {
		log.Printf("Error confirming user: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError, Body: err.Error()}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: "User confirmed successfully"}, nil
}
