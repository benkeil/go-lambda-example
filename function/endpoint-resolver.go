package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/service/sqs"
	"os"
)

var defaultResolver = endpoints.DefaultResolver()
var localstackUrl, found = os.LookupEnv("LOCALSTACK_URL")

const (
	signingRegion = "localstack"
	SqsPort = 4576
)

func LocalStackEndpointResolverFunc(service, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
	if found && service == sqs.ServiceName {
		fmt.Printf("Try to use localstack at %s for %s\n", localstackUrl, sqs.ServiceName)
		return endpointResolver(SqsPort)
	}
	return defaultResolver.EndpointFor(service, region, optFns...)
}

func endpointResolver(port int) (endpoints.ResolvedEndpoint, error) {
	return endpoints.ResolvedEndpoint{
		URL:           endpointUrl(port),
		SigningRegion: signingRegion,
	}, nil
}

func endpointUrl(port int) string {
	return fmt.Sprintf("%s:%d", localstackUrl, port)
}
