package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sqs"
	"os"
)

var defaultResolver = endpoints.DefaultResolver()
var localstackUrl, found = os.LookupEnv("LOCALSTACK_URL")

const (
	signingRegion = "localstack"
	SqsPort       = 4576
	IamPort       = 4593
)

func LocalStackEndpointResolverFunc(service, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
	if !found {
		return defaultResolver.EndpointFor(service, region, optFns...)
	}

	fmt.Printf("Try to use localstack at %s for %s\n", localstackUrl, sqs.ServiceName)

	// maybe use for all services a environment variable
	switch service {
	case sqs.ServiceName:
		return localstackEndpointResolver(SqsPort)
	case iam.ServiceName:
		return localstackEndpointResolver(IamPort)
	default:
		return defaultResolver.EndpointFor(service, region, optFns...)
	}
}

func localstackEndpointResolver(port int) (endpoints.ResolvedEndpoint, error) {
	return endpoints.ResolvedEndpoint{
		URL:           endpointUrl(port),
		SigningRegion: signingRegion,
	}, nil
}

func endpointUrl(port int) string {
	return fmt.Sprintf("%s:%d", localstackUrl, port)
}
