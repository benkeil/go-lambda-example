# go-lambda-example

## Testing

- Use [LocalStack](https://github.com/localstack/localstack) to mock most of AWS services.
- Implement custom endpoint resolver, e.g. based on a environment variable ([see this issue](https://github.com/awslabs/aws-sam-cli/issues/92)),
to use the SDK in a normal way but speak with LocalStack. [Here](https://docs.aws.amazon.com/sdk-for-go/api/aws/endpoints/) is an example for Go.
- 