package lib

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Birthday is the request body. the Username is optional because it is filled from the path parameter
type Birthday struct {
	Birthday string
	Username string `json:",omitempty"`
}

// NewResponse creates a new API Gateway Response. If marshalling fails, it returns with the error.
func NewResponse(status int, body interface{}) (events.APIGatewayProxyResponse, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(data),
	}, nil
}

// NewJsonResponse creates a new API Gateway Response that the APIG can transform into a proper response.
func NewJsonResponse(status int, message string) (events.APIGatewayProxyResponse, error) {
	ok := false

	if status > 199 && status < 400 {
		ok = true
	}

	return NewResponse(status, struct {
		Message string `json:"message"`
		Ok      bool   `json:"ok"`
	}{
		Message: message,
		Ok:      ok,
	})
}

// NewDynamoDbClient creates a new dynamodb client from the current session
func NewDynamoDbClient() (*dynamodb.DynamoDB, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	})
	if err != nil {
		return nil, err
	}

	return dynamodb.New(sess), nil
}
