package user

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
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

func CreateUser() {

}

func UpdateUser() {

}

func DeleteUser() error {

}
