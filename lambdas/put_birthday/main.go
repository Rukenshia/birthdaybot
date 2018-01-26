package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	log "github.com/sirupsen/logrus"
)

// PutBirthdayBody is the request body. the Username is optional because it is filled from the path parameter
type PutBirthdayBody struct {
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

// Handler handles the Amazon API Gateway Event
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Infof("Processing Lambda request %s", request.RequestContext.RequestID)
	log.Debugf("Request body: %s", request.Body)

	var body PutBirthdayBody
	if err := json.Unmarshal([]byte(request.Body), &body); err != nil {
		log.Errorf("JSON unmarshal failed: %v", err)
		return NewJsonResponse(400, "Invalid body. Expected JSON")
	}

	// Try to parse the Birthday
	log.Debugf("Attempting to parse Birthday")
	if _, err := time.Parse("2006-01-02", body.Birthday); err != nil {
		log.Errorf("Invalid Birthday: %v", err)
		return NewJsonResponse(422, "Invalid Birthday")
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	})
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	body.Username = request.PathParameters["Username"]

	dyndb := dynamodb.New(sess)

	item, err := dynamodbattribute.MarshalMap(body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	_, err = dyndb.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("DatabaseName")),
		Item:      item,
	})
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return NewJsonResponse(200, "Birthday added")
}

func main() {
	log.SetLevel(log.DebugLevel)
	lambda.Start(Handler)
}
