package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var (
	// DefaultHTTPGetAddress Default Address
	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"

	// ErrNoIP No IP found in response
	ErrNoIP = errors.New("no IP in HTTP response")

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("non 200 Response found")

	//QueueId the name of the queue
	QueueId = os.Getenv("QUEUE_ID")

	//QueueUrl the URL of the queue
	QueueUrl = os.Getenv("QUEUE_URL")
)

type HelloResponse struct {
	Name string
	IP string
	Message string
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("Request: %+v\n", request)

	name := request.QueryStringParameters["name"]
	if name == "" {
		name = "User"
	}

	resp, err := http.Get(DefaultHTTPGetAddress)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if resp.StatusCode != 200 {
		return events.APIGatewayProxyResponse{}, ErrNon200Response
	}

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if len(ip) == 0 {
		return events.APIGatewayProxyResponse{}, ErrNoIP
	}
	ipAddress := strings.TrimSpace(string(ip))

	response := HelloResponse{
		Name: name,
		IP: ipAddress,
		Message: fmt.Sprintf("Hello %s, your IP is %v", name, ipAddress),
	}
	body, err := json.Marshal(response)

	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(endpoints.EuCentral1RegionID),
		EndpointResolver: endpoints.ResolverFunc(LocalStackEndpointResolverFunc),
	})
	sqsService := sqs.New(awsSession)

	//queueUrl, err := sqsService.GetQueueUrl(&sqs.GetQueueUrlInput{
	//	QueueName: aws.String(QueueId),
	//})
	//if err != nil {
	//	return events.APIGatewayProxyResponse{}, err
	//}
	//fmt.Printf("URL for queue name %s -> %s\n", QueueId, aws.StringValue(queueUrl.QueueUrl))

	fmt.Printf("Sending message to queue %s\n", QueueId)
	sqsResponse, err := sqsService.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(string(body)),
		QueueUrl: &QueueUrl,
	})
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	fmt.Printf("SQS Response: %+v\n", sqsResponse)

	receivedMessage, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl: &QueueUrl,
		WaitTimeSeconds: aws.Int64(5),
	})
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	fmt.Printf("Received message: %+v\n", receivedMessage)

	return events.APIGatewayProxyResponse{
		Body: string(body),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(handler)
}
