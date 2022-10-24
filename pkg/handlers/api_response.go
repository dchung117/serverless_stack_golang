package handlers

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

func apiResponse(status int, body interface{}) (*events.APIGatewayProxyResponse, error) {
	// set up JSON response w/ status code
	response := events.APIGatewayProxyResponse{Headers: map[string]string{"Content-Type": "application/json"}}
	response.StatusCode = status

	// Marshal response body from Golang to JSON
	stringBody, _ := json.Marshal(body)
	response.Body = string(stringBody)

	return &response, nil
}
