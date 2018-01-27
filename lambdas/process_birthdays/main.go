package main

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	log "github.com/sirupsen/logrus"

	"github.com/Rukenshia/birthdaybot/lambdas/lib"
)

// Handler handles the Amazon API Gateway Event
func Handler(event events.CloudWatchEvent) (interface{}, error) {
	rlog := log.WithField("EventId", event.ID).WithField("DatabaseName", os.Getenv("DatabaseName"))

	rlog.Infof("Handling CloudWatch event")

	dyndb, err := lib.NewDynamoDbClient()
	if err != nil {
		log.Errorf("Creating DynDB Client failed: %v", err)
		return nil, err
	}

	rlog.Debugf("Scanning DynamoDB for all birthdays")
	results, err := dyndb.Scan(&dynamodb.ScanInput{
		TableName: aws.String(os.Getenv("DatabaseName")),
		ExpressionAttributeNames: map[string]*string{
			"#UN": aws.String("Username"),
			"#BD": aws.String("Birthday"),
		},
		ProjectionExpression: aws.String("#UN, #BD"),
	})
	if err != nil {
		rlog.Errorf("Scan failed: %v", err)
		return nil, err
	}

	rlog.Debugf("Scan returned %d result(s)", *results.Count)

	return nil, nil
}

func main() {
	log.SetLevel(log.DebugLevel)
	Handler(events.CloudWatchEvent{ID: "x"})
	lambda.Start(Handler)
}
