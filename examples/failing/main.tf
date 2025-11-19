# Example of bad practices - hardcoded regions and partitions

terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}
#
# # Hardcoded region in provider (will trigger aws_provider_hardcoded_region rule)
# provider "aws" {
#   region = "us-east-1"
# }
#
# # Provider with hardcoded region in assume_role
# provider "aws" {
#   alias = "assume_role_hardcoded"
#   assume_role {
#     role_arn = "arn:aws:iam:us-west-2:123456789012:role/terraform-role"
#   }
# }
#
# # IAM role policy with hardcoded region (will trigger aws_iam_role_policy_hardcoded_region rule)
# resource "aws_iam_role_policy" "example_region" {
#   name = "example-policy-region"
#   role = "example-role"
#
#   policy = "{\"Version\": \"2012-10-17\", \"Statement\": [{\"Action\": [\"s3:GetObject\"], \"Effect\": \"Allow\", \"Resource\": \"arn:aws:s3:::my-bucket/us-east-1/*\"}]}"
# }
#
# # IAM role policy with hardcoded partition (will trigger aws_iam_role_policy_hardcoded_partition rule)
# resource "aws_iam_role_policy" "example_partition" {
#   name = "example-policy-partition"
#   role = "example-role"
#
#   policy = "{\"Version\": \"2012-10-17\", \"Statement\": [{\"Action\": [\"lambda:InvokeFunction\"], \"Effect\": \"Allow\", \"Resource\": \"arn:aws:lambda:us-west-2:123456789012:function:my-function\"}]}"
# }
#
# # IAM policy with hardcoded region (will trigger aws_iam_policy_hardcoded_region rule)
# resource "aws_iam_policy" "example_region" {
#   name = "example-policy-region"
#
#   policy = "{\"Version\": \"2012-10-17\", \"Statement\": [{\"Action\": [\"dynamodb:GetItem\"], \"Effect\": \"Allow\", \"Resource\": \"arn:aws:dynamodb:eu-west-1:123456789012:table/my-table\"}]}"
# }
#
# # IAM policy with hardcoded partition (will trigger aws_iam_policy_hardcoded_partition rule)
# resource "aws_iam_policy" "example_partition" {
#   name = "example-policy-partition"
#
#   policy = "{\"Version\": \"2012-10-17\", \"Statement\": [{\"Action\": [\"sqs:SendMessage\"], \"Effect\": \"Allow\", \"Resource\": \"arn:aws-us-gov:sqs:us-gov-west-1:123456789012:my-queue\"}]}"
# }
#
# # Multiple violations in one policy
# resource "aws_iam_role_policy" "multiple_violations" {
#   name = "multiple-violations"
#   role = "example-role"
#
#   policy = "{\"Version\": \"2012-10-17\", \"Statement\": [{\"Action\": [\"s3:GetObject\"], \"Effect\": \"Allow\", \"Resource\": [\"arn:aws:s3:::bucket1/ap-southeast-1/*\", \"arn:aws-cn:s3:::bucket2/cn-north-1/*\"]}]}"
# }
#
#
# # Lambda permission with hardcoded ARN (will trigger aws_arn_hardcoded rule)
# resource "aws_lambda_permission" "example" {
#   statement_id  = "AllowS3Invoke"
#   action        = "lambda:InvokeFunction"
#   function_name = "my-function"
#   principal     = "s3.amazonaws.com"
#   source_arn    = "arn:aws:s3:us-east-1:123456789012:bucket/my-bucket"
# }
#
# # SNS subscription with hardcoded ARN (will trigger aws_arn_hardcoded rule)
# resource "aws_sns_topic_subscription" "example" {
#   topic_arn = "arn:aws:sns:eu-west-1:123456789012:my-topic"
#   protocol  = "email"
#   endpoint  = "example@example.com"
# }

# KMS grant with hardcoded ARN (will trigger aws_arn_hardcoded rule)
resource "aws_kms_grant" "example" {
  name              = "my-grant"
  key_id            = "arn:aws:kms:ap-southeast-1:123456789012:key/12345678-1234-1234-1234-123456789012"
  grantee_principal = "arn:aws:iam::123456789012:role/my-role"
}

# CloudWatch event target with hardcoded ARN (will trigger aws_arn_hardcoded rule)
resource "aws_cloudwatch_event_target" "example" {
  rule = "my-rule"
  arn  = "arn:aws-cn:lambda:cn-north-1:123456789012:function:my-function"
}
