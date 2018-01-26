BINARY_NAME=main
DEPLOYMENT_BUCKET=birthdaybot-artifacts

all: clean build zip

upload:
	aws s3 cp lambda.zip s3://${DEPLOYMENT_BUCKET}/put_birthday/lambda.zip

clean:
	rm -f ${BINARY_NAME} lambda.zip

build:
	go build -o ${BINARY_NAME}

zip: build
	zip lambda.zip ${BINARY_NAME}
