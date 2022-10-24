package main

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	"github.com/dchung117/serverless_stack_golang/pkg/handlers"
)

// create dynaClient
var (
	dynaClient dynamodbiface.DynamoDBAPI
)

const tableName = "LambdaInGoUser"

func main() {
	// get AWS REGION from environment
	region := os.Getenv("AWS_REGION")

	// create AWS session
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return
	}

	// create dynamodb instance w/ AWS session
	dynaClient = dynamodb.New(awsSession)

	// start aws lambda handler
	lambda.Start(handler)
}

// handler to route client requests to appropriate handle function
func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return handlers.GetUser(req, tableName, dynaClient)
	case "POST":
		return handlers.CreateUser(req, tableName, dynaClient)
	case "PUT":
		return handlers.UpdateUser(req, tableName, dynaClient)
	case "DELETE":
		return handlers.DeleteUser(req, tableName, dynaClient)
	default:
		return handlers.UnhandledMethod()
	}

}
