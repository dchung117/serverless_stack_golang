package user

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/dchung117/serverless_stack_golang/pkg/validators"
)

// define user struct
type User struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// define custom errors
var (
	ErrorFailedToFetchRecord     = "failed to fetch record."
	ErrorFailedtoUnmarshalRecord = "failed to unmarshal record."
	ErrorInvalidUserData         = "invalid user data."
	ErrorInvalidEmail            = "invalid email."
	ErrorFailedtoMarshalItem     = "failed to marshal item."
	ErrorCouldNotDeleteItem      = "failed to delete item."
	ErrorCouldNotDynamoPutItem   = "failed to dynamo put item."
	ErrorUserAlreadyExists       = "user already exists."
	ErrorUserDoesNotExist        = "user does not exist."
)

func FetchUser(email string, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*User, error) {
	// set up query to retrieve user data
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}

	// get user data
	result, err := dynaClient.GetItem(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	// create a new user, unmarshal result into item
	item := new(User)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, errors.New(ErrorFailedtoUnmarshalRecord)
	}

	// return new user struct
	return item, nil
}

func FetchUsers(tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*[]User, error) {
	// set up query for getting all users
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	// find all users
	result, err := dynaClient.Scan(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	// create slice of users, unmarshal results into slice for JSON
	item := new([]User)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, item)
	if err != nil {
		return nil, errors.New(ErrorFailedtoUnmarshalRecord)
	}
	return item, nil
}

func CreateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*User, error) {
	// create new user variable
	var user User

	// unmarshal user info from request body
	if err := json.Unmarshal([]byte(req.Body), &user); err != nil {
		return nil, errors.New(ErrorFailedtoUnmarshalRecord)
	}

	// validate the email
	if !validators.IsEmailValid(user.Email) {
		return nil, errors.New(ErrorInvalidEmail)
	}

	// check if user already exists
	currentUser, _ := FetchUser(user.Email, tableName, dynaClient)
	if (currentUser != nil) && (len(currentUser.Email) != 0) {
		return nil, errors.New(ErrorUserAlreadyExists)
	}

	// marshal the result
	avMap, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return nil, errors.New(ErrorFailedtoMarshalItem)
	}

	// update database
	input := &dynamodb.PutItemInput{
		Item:      avMap,
		TableName: aws.String(tableName),
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}

	return &user, nil
}

func UpdateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*User, error) {
	// create a new user
	var user User

	// unmarshal the request body
	if err := json.Unmarshal([]byte(req.Body), &user); err != nil {
		return nil, errors.New(ErrorFailedtoUnmarshalRecord)
	}

	// validate the email
	if !validators.IsEmailValid(user.Email) {
		return nil, errors.New(ErrorInvalidEmail)
	}

	// check if user exists
	currentUser, _ := FetchUser(user.Email, tableName, dynaClient)
	if currentUser != nil && len(currentUser.Email) == 0 {
		return nil, errors.New(ErrorUserDoesNotExist)
	}

	// marshal the result
	avMap, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return nil, errors.New(ErrorFailedtoMarshalItem)
	}

	// update database
	input := &dynamodb.PutItemInput{
		Item:      avMap,
		TableName: aws.String(tableName),
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}

	return &user, nil
}

func DeleteUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) error {
	// get email from query parameters
	email := req.QueryStringParameters["email"]

	// delete user w/ matching email
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}

	_, err := dynaClient.DeleteItem(input)
	if err != nil {
		return errors.New(ErrorCouldNotDeleteItem)
	}

	return nil
}
