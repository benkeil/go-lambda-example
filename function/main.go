package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	// DefaultHTTPGetAddress Default Address
	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"

	// ErrNoIP No IP found in response
	ErrNoIP = errors.New("no IP in HTTP response")

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("non 200 Response found")
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
	return events.APIGatewayProxyResponse{
		Body: string(body),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
