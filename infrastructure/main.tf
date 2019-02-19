provider "aws" {
  region  = "${var.aws_region}"
  version = "~> 1.59"
  endpoints {
    sqs = "${local.aws_endpoint_sqs}"
    s3 = "${local.aws_endpoint_s3}"
    lambda = "${local.aws_endpoint_lambda}"
  }
  max_retries = 1
}

data "aws_caller_identity" "current" {}