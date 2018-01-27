package main

import (
	"fmt"
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
func Handler(event events.CloudWatchEvent) (interface{}, error) {
	rlog := log.WithField("EventId", event.ID).WithField("DatabaseName", os.Getenv("DatabaseName"))
	rlog.Infof("Handling CloudWatch event")

	url, err := lib.DecryptEnvVar("", "WebhookUrl")
	if err != nil {
		rlog.Errorf("Decrypting WebhookUrl failed: %v", err)
		return nil, err
	}

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

	var birthdays []lib.Birthday
	if err := dynamodbattribute.UnmarshalListOfMaps(results.Items, &birthdays); err != nil {
		return nil, err
	}

	today := time.Now()

	for _, birthday := range birthdays {
		date, err := time.Parse("2006-01-02", birthday.Birthday)

		if err != nil {
			rlog.Errorf("Invalid database birthday for %s: %s", birthday.Username, birthday.Birthday)
			return nil, err
		}

		rlog.Debugf("%d == %d && %d == %d", date.Month(), today.Month(), date.Day(), today.Day())

		if date.Month() == today.Month() && date.Day() == today.Day() {
			if err := SendRocketchatMessage(url, fmt.Sprintf("@%s", birthday.Username), "whoooo!"); err != nil {
				rlog.WithField("Username", birthday.Username).Warnf("Could not send Rocket.Chat notification: %v", err)
			}
		}
	}

	return nil, nil
}

func main() {
	log.SetLevel(log.DebugLevel)
	lambda.Start(Handler)
}
