resource "aws_sqs_queue" "hello" {
  name                       = "hello"
  delay_seconds              = 1
  max_message_size           = 2048
  message_retention_seconds  = 86400
  receive_wait_time_seconds  = 10
  visibility_timeout_seconds = 300
  redrive_policy             = "{\"deadLetterTargetArn\":\"${aws_sqs_queue.hello_dl.arn}\",\"maxReceiveCount\":4}"
}

resource "aws_sqs_queue" "hello_dl" {
  name = "hello-dlq"
}

resource "aws_sqs_queue_policy" "hello_queue" {
  count = "${var.use_localstack ? 0 : 1}"
  queue_url = "${aws_sqs_queue.hello.id}"

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Id": "sqspolicy",
  "Statement": [
    {
      "Sid": "sns_sqs",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "sqs:SendMessage",
      "Resource": "${aws_sqs_queue.hello.arn}",
      "Condition": {
        "ArnEquals": {
          "aws:SourceArn": "${aws_lambda_function.hello.arn}"
        }
      }
    }
  ]
}
POLICY
}
