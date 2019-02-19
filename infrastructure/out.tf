output "queue_id" {
  value = "${aws_sqs_queue.hello.id}"
}

output "lambda_arn" {
  value = "${aws_lambda_function.hello.*.arn}"
}

resource "local_file" "output" {
  content  = <<EOF
{
  "queue_id": "${aws_sqs_queue.hello.id}",
  "lambda_arn": "${join("", aws_lambda_function.hello.*.arn)}"
}
EOF
  filename = "${path.module}/../build/output.json"
}