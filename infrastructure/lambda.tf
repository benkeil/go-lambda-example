resource "aws_lambda_function" "hello" {
  #count = "${var.use_localstack ? 0 : 1}"
  filename         = "lambda.zip"
  function_name    = "hello"
  role             = "${aws_iam_role.lambda.arn}"
  handler          = "main"
  #source_code_hash = "${base64sha256(file(../lambda.zip"))}"
  runtime          = "go1.x"

  environment {
    variables = {
      LOCALSTACK_URL = "http://192.168.178.21"
      QUEUE_ID = "hello"
      QUEUE_URL = "${aws_sqs_queue.hello.id}"
    }
  }

  depends_on = [
    "data.archive_file.lambda_zip",
    "aws_iam_role.lambda",
  ]
}

data "archive_file" "lambda_zip" {
  #count = "${var.use_localstack ? 0 : 1}"
  type                    = "zip"
  source_content          = "hello lambda"
  source_content_filename = "main"
  output_path             = "lambda.zip"
}

resource "aws_iam_role" "lambda" {
  #count = "${var.use_localstack ? 0 : 1}"
  name = "hello-lambda"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}