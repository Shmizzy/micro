package main

import (
	"service/handlers"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	lambda.Start(handlers.ConfirmHandler)
}
