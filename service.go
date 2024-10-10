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

	createUser := awslambda.NewFunction(stack, jsii.String("confirmFunction"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("main"),
		Code:    awslambda.AssetCode_FromAsset(jsii.String("cmd/create_user/createUser.zip"), nil),
	})

	getUser := awslambda.NewFunction(stack, jsii.String("getUserFunction"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("main"),
		Code:    awslambda.AssetCode_FromAsset(jsii.String("cmd/get_user/getUser.zip"), nil),
	})

	updateUser := awslambda.NewFunction(stack, jsii.String("updateUserFunction"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("main"),
		Code:    awslambda.AssetCode_FromAsset(jsii.String("cmd/update_user/updateUser.zip"), nil),
	})

	deleteUser := awslambda.NewFunction(stack, jsii.String("deleteUserFunction"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("main"),
		Code:    awslambda.AssetCode_FromAsset(jsii.String("cmd/delete_user/deleteUser.zip"), nil),
	})

	createServicerProfile := awslambda.NewFunction(stack, jsii.String("createServicerProfileFunction"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("main"),
		Code:    awslambda.AssetCode_FromAsset(jsii.String("cmd/create_servicer_profile/createServicerProfile.zip"), nil),
	})

	updateServicerProfile := awslambda.NewFunction(stack, jsii.String("updateServicerProfileFunction"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("main"),
		Code:    awslambda.AssetCode_FromAsset(jsii.String("cmd/update_servicer_profile/updateServicerProfile.zip"), nil),
	})

	getServicerProfile := awslambda.NewFunction(stack, jsii.String("getServicerProfileFunction"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("main"),
		Code:    awslambda.AssetCode_FromAsset(jsii.String("cmd/get_servicer_profile/getServicerProfile.zip"), nil),
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

	userTable := awsdynamodb.NewTableV2(stack, jsii.String("UserTable"), &awsdynamodb.TablePropsV2{
		TableName: jsii.String("usersTable"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("UserId"),
			Type: awsdynamodb.AttributeType_STRING,
		},
	})

	customerProfilesTable := awsdynamodb.NewTableV2(stack, jsii.String("CustomerProfileTable"), &awsdynamodb.TablePropsV2{
		TableName: jsii.String("customerProfilesTable"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("UserId"),
			Type: awsdynamodb.AttributeType_STRING,
		},
	})

	servicerProfilesTable := awsdynamodb.NewTableV2(stack, jsii.String("ServicerProfileTable"), &awsdynamodb.TablePropsV2{
		TableName: jsii.String("servicerProfilesTable"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("UserId"),
			Type: awsdynamodb.AttributeType_STRING,
		},
	})
	postConfirmationFunction := awslambda.NewFunction(stack, jsii.String("PostConfirmationFunction"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("main"),
		Code:    awslambda.AssetCode_FromAsset(jsii.String("cmd/post_confirmation/postConfirmation.zip"), nil),
		Environment: &map[string]*string{
			"USERS_TABLE":             userTable.TableName(),
			"CUSTOMER_PROFILES_TABLE": customerProfilesTable.TableName(),
		},
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
			Email: jsii.Bool(true),
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
			EmailBody:    jsii.String("Hello {name}, your verification code is {####}"),
			SmsMessage:   jsii.String("Hello {name}, your verification code is {####}"),
			EmailStyle:   awscognito.VerificationEmailStyle_CODE,
		},
	})

	auth := awsapigateway.NewCognitoUserPoolsAuthorizer(stack, jsii.String("UserPoolAuthorizer"), &awsapigateway.CognitoUserPoolsAuthorizerProps{
		CognitoUserPools: &[]awscognito.IUserPool{
			userPool,
		},
	})

	userTable.GrantWriteData(postConfirmationFunction)
	customerProfilesTable.GrantWriteData(postConfirmationFunction)
	userTable.GrantReadWriteData(createUser)
	userTable.GrantReadData(getUser)
	userTable.GrantReadWriteData(updateUser)
	userTable.GrantReadWriteData(deleteUser)
	customerProfilesTable.GrantReadWriteData(createUser)
	customerProfilesTable.GrantReadData(getUser)
	customerProfilesTable.GrantReadWriteData(updateUser)
	servicerProfilesTable.GrantReadWriteData(createServicerProfile)
	servicerProfilesTable.GrantReadWriteData(updateServicerProfile)
	servicerProfilesTable.GrantReadData(getServicerProfile)

	createUser.AddEnvironment(jsii.String("USERS_TABLE"), userTable.TableName(), nil)
	createUser.AddEnvironment(jsii.String("CUSTOMER_PROFILES_TABLE"), customerProfilesTable.TableName(), nil)
	getUser.AddEnvironment(jsii.String("USERS_TABLE"), userTable.TableName(), nil)
	getUser.AddEnvironment(jsii.String("CUSTOMER_PROFILES_TABLE"), customerProfilesTable.TableName(), nil)
	updateUser.AddEnvironment(jsii.String("USERS_TABLE"), userTable.TableName(), nil)
	updateUser.AddEnvironment(jsii.String("CUSTOMER_PROFILES_TABLE"), customerProfilesTable.TableName(), nil)
	deleteUser.AddEnvironment(jsii.String("USERS_TABLE"), userTable.TableName(), nil)
	deleteUser.AddEnvironment(jsii.String("CUSTOMER_PROFILES_TABLE"), customerProfilesTable.TableName(), nil)
	createServicerProfile.AddEnvironment(jsii.String("SERVICER_PROFILES_TABLE"), servicerProfilesTable.TableName(), nil)
	updateServicerProfile.AddEnvironment(jsii.String("SERVICER_PROFILES_TABLE"), servicerProfilesTable.TableName(), nil)
	getServicerProfile.AddEnvironment(jsii.String("SERVICER_PROFILES_TABLE"), servicerProfilesTable.TableName(), nil)

	usersResource := api.Root().AddResource(jsii.String("users"), nil)
	servicerProfilesResource := api.Root().AddResource(jsii.String("servicer-profiles"), nil)

	userPool.AddTrigger(awscognito.UserPoolOperation_POST_CONFIRMATION(), postConfirmationFunction, awscognito.LambdaVersion_V1_0)

	// route -- users
	usersResource.AddMethod(jsii.String("POST"), awsapigateway.NewLambdaIntegration(createUser, nil), &awsapigateway.MethodOptions{
		AuthorizationType: awsapigateway.AuthorizationType_COGNITO,
		Authorizer:        auth,
	})

	// route -- users/{userId}
	userResource := usersResource.AddResource(jsii.String("{userId}"), nil)

	userResource.AddMethod(jsii.String("GET"), awsapigateway.NewLambdaIntegration(getUser, nil), &awsapigateway.MethodOptions{
		AuthorizationType: awsapigateway.AuthorizationType_COGNITO,
		Authorizer:        auth,
	})

	userResource.AddMethod(jsii.String("PUT"), awsapigateway.NewLambdaIntegration(updateUser, nil), &awsapigateway.MethodOptions{
		AuthorizationType: awsapigateway.AuthorizationType_COGNITO,
		Authorizer:        auth,
	})

	userResource.AddMethod(jsii.String("DELETE"), awsapigateway.NewLambdaIntegration(deleteUser, nil), &awsapigateway.MethodOptions{
		AuthorizationType: awsapigateway.AuthorizationType_COGNITO,
		Authorizer:        auth,
	})

	// route -- servicer-profiles
	servicerProfilesResource.AddMethod(jsii.String("POST"), awsapigateway.NewLambdaIntegration(createServicerProfile, nil), &awsapigateway.MethodOptions{
		AuthorizationType: awsapigateway.AuthorizationType_COGNITO,
		Authorizer:        auth,
	})

	// route -- servicer-profiles/{userId}
	servicerProfileResource := servicerProfilesResource.AddResource(jsii.String("{userId}"), nil)

	servicerProfileResource.AddMethod(jsii.String("GET"), awsapigateway.NewLambdaIntegration(getServicerProfile, nil), &awsapigateway.MethodOptions{
		AuthorizationType: awsapigateway.AuthorizationType_COGNITO,
		Authorizer:        auth,
	})

	servicerProfileResource.AddMethod(jsii.String("PUT"), awsapigateway.NewLambdaIntegration(updateServicerProfile, nil), &awsapigateway.MethodOptions{
		AuthorizationType: awsapigateway.AuthorizationType_COGNITO,
		Authorizer:        auth,
	})

	awscdk.NewCfnOutput(stack, jsii.String("UserPoolId"), &awscdk.CfnOutputProps{
		Value:       userPool.UserPoolId(),
		Description: jsii.String("The ID of the Cognito User Pool"),
		ExportName:  jsii.String("UserPoolId"),
	})

	awscdk.NewCfnOutput(stack, jsii.String("ApiUrl"), &awscdk.CfnOutputProps{
		Value:       api.Url(),
		Description: jsii.String("The URL of the API Gateway"),
		ExportName:  jsii.String("ApiUrl"),
	})

	/* userPool.AddTrigger(awscognito.UserPoolOperation_PRE_SIGN_UP(), registerFunction, awscognito.LambdaVersion_V1_0)
	userPool.AddTrigger(awscognito.UserPoolOperation_PRE_AUTHENTICATION(), loginFunction, awscognito.LambdaVersion_V1_0)
	*/

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
