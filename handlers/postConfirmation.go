package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"service/models"

	"github.com/aws/aws-lambda-go/events"
)

func PostConfirmationHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Recieved data -> %s", request.Body)
	var user models.User
	err := json.Unmarshal([]byte(request.Body), &user)
	if err != nil {
		log.Printf("Error unmarshalling data -> %v", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest, Body: "Invalid request body"}, nil
	}

	err = services.CreateUserEntry(user)
	if err != nil {
		log.Printf("Error creating user entry -> %v", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError, Body: err.Error()}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: "User created successfully"}, nil
}
