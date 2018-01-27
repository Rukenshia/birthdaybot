package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	log "github.com/sirupsen/logrus"

	"github.com/Rukenshia/birthdaybot/lambdas/lib"
)

// Handler handles the Amazon API Gateway Event
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Infof("Processing Lambda request %s", request.RequestContext.RequestID)
	log.Debugf("Request body: %s", request.Body)

	var body lib.Birthday
	if err := json.Unmarshal([]byte(request.Body), &body); err != nil {
		log.Errorf("JSON unmarshal failed: %v", err)
		return lib.NewJsonResponse(400, "Invalid body. Expected JSON")
	}

	// Try to parse the Birthday
	log.Debugf("Attempting to parse Birthday")
	if _, err := time.Parse("2006-01-02", body.Birthday); err != nil {
		log.Errorf("Invalid Birthday: %v", err)
		return lib.NewJsonResponse(422, "Invalid Birthday")
	}

	dyndb, err := lib.NewDynamoDbClient()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	body.Username = request.PathParameters["Username"]
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

	return lib.NewJsonResponse(200, "Birthday added")
}

func main() {
	log.SetLevel(log.DebugLevel)
	lambda.Start(Handler)
}
