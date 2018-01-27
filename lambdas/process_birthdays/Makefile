BINARY_NAME=main
DEPLOYMENT_BUCKET=birthdaybot-artifacts
LAMBDA_PATH=process_birthdays/lambda.zip
BASE_URL=https://birthday.in.fkn.space

all: clean build zip

artifact: zip
	aws s3 cp lambda.zip s3://$(DEPLOYMENT_BUCKET)/$(LAMBDA_PATH)

upload: artifact
	$(eval function_name := `aws lambda list-functions | jq -r '.Functions[] | select(.FunctionName | contains("ProcessBirthdays")).FunctionName'`)

	aws lambda update-function-code \
		--function-name $(function_name) \
		--s3-bucket $(DEPLOYMENT_BUCKET) \
		--s3-key $(LAMBDA_PATH) \
		--publish

clean:
	rm -f $(BINARY_NAME) lambda.zip

build:
	GOOS=linux go build -o $(BINARY_NAME)

zip: build
	zip lambda.zip $(BINARY_NAME)
