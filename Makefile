# Referenced as Handler in template.yaml
OUTPUT = main
PACKAGED_TEMPLATE = packaged.yaml
S3_BUCKET := benkeil-cloudformation
STACK_NAME := lambda-hello
TEMPLATE = template.yaml
FUNCTION_NAME = HelloFunction
LDFLAGS :=

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	rm -f $(OUTPUT) $(PACKAGED_TEMPLATE)

.PHONY: install
install:
	go get ./...

.PHONY: main
main: ./function/main.go
	go build -ldflags="$(LDFLAGS)" -o $(OUTPUT) ./function/main.go

# compile the code to run in Lambda (local or real)
.PHONY: lambda
lambda:
	GOOS=linux GOARCH=amd64 $(MAKE) main LDFLAGS="$(LDFLAGS)"

.PHONY: build
build: clean lambda

.PHONY: api
start-api: build
	sam local start-api

.PHONY:
start-lambda: build
	sam local start-lambda

.PHONY: package
package: build
	sam package --template-file $(TEMPLATE) --s3-bucket $(S3_BUCKET) --output-template-file $(PACKAGED_TEMPLATE)

.PHONY: upx
upx: package
	upx --brute main

.PHONY: deploy
deploy: LDFLAGS = -s -w
deploy: upx
	sam deploy --stack-name $(STACK_NAME) --template-file $(PACKAGED_TEMPLATE) --capabilities CAPABILITY_IAM

.PHONY: invoke
invoke:
	aws lambda invoke --function-name $(FUNCTION_NAME) --endpoint-url http://127.0.0.1:3001 --payload file://test/events/aws-proxy.json test/response.json