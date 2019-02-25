variable "aws_region" {
  description = "The AWS region"
  type = "string"
  default = "eu-central-1"
}

variable "use_localstack" {
  description = "Should we use LocalStack?"
  default = true
}

variable "localstack_url" {
  description = "The URL of localstack"
  type = "string"
  default = "localhost"
}

# For endpoints see https://docs.aws.amazon.com/de_de/general/latest/gr/rande.html
locals {
  aws_endpoint_sqs = "${var.use_localstack ? "http://${var.localstack_url}:4576" : "https://sqs.${var.aws_region}.amazonaws.com"}"
  aws_endpoint_s3 = "${var.use_localstack ? "http://${var.localstack_url}:4572" : "https://s3.${var.aws_region}.amazonaws.com"}"
  aws_endpoint_lambda = "${var.use_localstack ? "http://${var.localstack_url}:4574" : "https://lambda.${var.aws_region}.amazonaws.com"}"
  aws_endpoint_iam = "${var.use_localstack ? "http://${var.localstack_url}:4593" : "https://iam.${var.aws_region}.amazonaws.com"}"
}