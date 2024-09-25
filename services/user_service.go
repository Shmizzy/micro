package services

import (
	"errors"
	"fmt"
	"log"
	"os"
	"service/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/joho/godotenv"
)

var (
	cognitoClient *cognitoidentityprovider.CognitoIdentityProvider
	dynamoClient  *dynamodb.DynamoDB
	userPoolId    string
	clientId      string
	dynamoTable   string
)

func init() {
	godotenv.Load("../.env")

	userPoolId = os.Getenv("USER_POOL_ID")
	clientId = os.Getenv("CLIENT_ID")
	dynamoTable = os.Getenv("DYNAMO_TABLE")

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	}))
	cognitoClient = cognitoidentityprovider.New(sess)
	dynamoClient = dynamodb.New(sess)
}

func RegisterUser(user models.User) error {
	log.Printf("Registering user: %v\n", user)
	signUp := &cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(clientId),
		Username: aws.String(user.Username),
		Password: aws.String(user.Password),
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(user.Email),
			},
			{
				Name:  aws.String("phone_number"),
				Value: aws.String(user.PhoneNumber),
			},
		},
	}

	res, err := cognitoClient.SignUp(signUp)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == cognitoidentityprovider.ErrCodeInvalidLambdaResponseException {
			return fmt.Errorf("InvalidLambdaResponseException: %v", aerr.Message())
		}
		return fmt.Errorf("failed to sign up user: %v", err)
	}
	fmt.Printf("SignUp response: %v\n", res)

	userInfo := map[string]interface{}{
		"email": user.Email,
		"uid":   user.Username,
	}

	av, err := dynamodbattribute.MarshalMap(userInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal user info: %v", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(dynamoTable),
		Item:      av,
	}
	_, err = dynamoClient.PutItem(input)
	if err != nil {
		return fmt.Errorf("failed to put user info in dynamodb: %v", err)
	}

	return nil
}

func ConfirmUser(username, confirmationCode string) error {
    input := &cognitoidentityprovider.ConfirmSignUpInput{
        ClientId:         aws.String(clientId),
        Username:         aws.String(username),
        ConfirmationCode: aws.String(confirmationCode),
    }

    _, err := cognitoClient.ConfirmSignUp(input)
    if err != nil {
        return fmt.Errorf("failed to confirm user: %v", err)
    }

    return nil
}

func LoginUser(user models.Credentials) (string, error) {

	authenticate := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(user.Username),
			"PASSWORD": aws.String(user.Password),
		},
		ClientId: aws.String(clientId),
	}

	resp, err := cognitoClient.InitiateAuth(authenticate)
	if err != nil {
		return "", fmt.Errorf("failed to authenticate user: %v", err)
	}

	if resp.AuthenticationResult == nil {
		return "", errors.New("authentication failed")
	}

	return *resp.AuthenticationResult.IdToken, nil

}
