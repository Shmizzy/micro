package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"

	// "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type ServiceStackProps struct {
	awscdk.StackProps
}

func NewUserServiceStack(scope constructs.Construct, id string, props *ServiceStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	registerFunction := awslambda.NewFunction(stack, jsii.String("registerFunction"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("main"),
		Code:    awslambda.AssetCode_FromAsset(jsii.String("cmd/register/register.zip"), nil),
	})

	loginFunction := awslambda.NewFunction(stack, jsii.String("loginFunction"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("main"),
		Code:    awslambda.AssetCode_FromAsset(jsii.String("cmd/login/login.zip"), nil),
	})

	confirmFunction := awslambda.NewFunction(stack, jsii.String("confirmFunction"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("main"),
		Code:    awslambda.AssetCode_FromAsset(jsii.String("cmd/confirm/confirm.zip"), nil),
	})

	userPool := awscognito.NewUserPool(stack, jsii.String("UserPool"), &awscognito.UserPoolProps{
		UserPoolName:        jsii.String("yardex-user-pool"),
		SelfSignUpEnabled:   jsii.Bool(true),
		SignInCaseSensitive: jsii.Bool(false),
		SignInAliases: &awscognito.SignInAliases{
			Username: jsii.Bool(true),
			Email:    jsii.Bool(true),
		},
		StandardAttributes: &awscognito.StandardAttributes{
			PhoneNumber: &awscognito.StandardAttribute{
				Required: jsii.Bool(true),
				Mutable:  jsii.Bool(true),
			},
			Email: &awscognito.StandardAttribute{
				Required: jsii.Bool(true),
				Mutable:  jsii.Bool(true),
			},
		},
		AutoVerify: &awscognito.AutoVerifiedAttrs{
			Email: jsii.Bool(false),
			Phone: jsii.Bool(true),
		},
		MfaSecondFactor: &awscognito.MfaSecondFactor{
			Sms: jsii.Bool(true),
			Otp: jsii.Bool(false),
		},
		Mfa:        awscognito.Mfa_OPTIONAL,
		MfaMessage: jsii.String("Your verification code for yardex is {####}"),
		PasswordPolicy: &awscognito.PasswordPolicy{
			MinLength:        jsii.Number(8),
			RequireLowercase: jsii.Bool(true),
			RequireUppercase: jsii.Bool(true),
			RequireDigits:    jsii.Bool(true),
			RequireSymbols:   jsii.Bool(false),
		},
		AccountRecovery: awscognito.AccountRecovery_PHONE_AND_EMAIL,
		UserVerification: &awscognito.UserVerificationConfig{
			EmailSubject: jsii.String("Verify your email for yardex"),
			EmailBody:    jsii.String("Hello {username}, your verification code is {####}"),
			SmsMessage:   jsii.String("Hello {username}, your verification code is {####}"),
			EmailStyle:   awscognito.VerificationEmailStyle_CODE,
		},
	})

	table := awsdynamodb.NewTableV2(stack, jsii.String("UserTable"), &awsdynamodb.TablePropsV2{
		TableName: jsii.String("userDB"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("uid"),
			Type: awsdynamodb.AttributeType_STRING,
		},
	})

	api := awsapigateway.NewRestApi(stack, jsii.String("UserApi"), &awsapigateway.RestApiProps{
		RestApiName: jsii.String("User Service API"),
		Description: jsii.String("API Gateway for User Service"),
		DefaultCorsPreflightOptions: &awsapigateway.CorsOptions{
			AllowOrigins: awsapigateway.Cors_ALL_ORIGINS(),
			AllowMethods: awsapigateway.Cors_ALL_METHODS(),
			AllowHeaders: jsii.Strings("Content-Type", "Authorization"),
		},
	})

	registerResource := api.Root().AddResource(jsii.String("register"), nil)
	loginResource := api.Root().AddResource(jsii.String("login"), nil)
	confirmResource := api.Root().AddResource(jsii.String("confirm"), nil)

	auth := awsapigateway.NewCognitoUserPoolsAuthorizer(stack, jsii.String("UserPoolAuthorizer"), &awsapigateway.CognitoUserPoolsAuthorizerProps{
		CognitoUserPools: &[]awscognito.IUserPool{
			userPool,
		},
	})

	registerIntegration := awsapigateway.NewLambdaIntegration(registerFunction, &awsapigateway.LambdaIntegrationOptions{
		RequestTemplates: &map[string]*string{
			"application/json": jsii.String(`{"statusCode": "200"}`),
		},
	})

	loginIntegration := awsapigateway.NewLambdaIntegration(loginFunction, &awsapigateway.LambdaIntegrationOptions{
		RequestTemplates: &map[string]*string{
			"application/json": jsii.String(`{"statusCode": "200"}`),
		},
	})

	confirmIntegration := awsapigateway.NewLambdaIntegration(confirmFunction, &awsapigateway.LambdaIntegrationOptions{
		RequestTemplates: &map[string]*string{
			"application/json": jsii.String(`{"statusCode": "200"}`),
		},
	})

	registerResource.AddMethod(jsii.String("POST"), registerIntegration, &awsapigateway.MethodOptions{
		AuthorizationType: awsapigateway.AuthorizationType_COGNITO,
		Authorizer:        auth,
	})
	loginResource.AddMethod(jsii.String("POST"), loginIntegration, &awsapigateway.MethodOptions{
		AuthorizationType: awsapigateway.AuthorizationType_COGNITO,
		Authorizer:        auth,
	})
	confirmResource.AddMethod(jsii.String("POST"), confirmIntegration, &awsapigateway.MethodOptions{
		AuthorizationType: awsapigateway.AuthorizationType_COGNITO,
		Authorizer:        auth,
	})

	/* userPool.AddTrigger(awscognito.UserPoolOperation_PRE_SIGN_UP(), registerFunction, awscognito.LambdaVersion_V1_0)
	userPool.AddTrigger(awscognito.UserPoolOperation_PRE_AUTHENTICATION(), loginFunction, awscognito.LambdaVersion_V1_0)
	*/
	table.GrantReadWriteData(registerFunction)
	table.GrantReadWriteData(loginFunction)
	table.GrantReadWriteData(confirmFunction)

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewUserServiceStack(app, "UserServiceStack", &ServiceStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {

	return nil

}
