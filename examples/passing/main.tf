# Example of good practices - no hardcoded regions or partitions

terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }
}

# Use variables for region configuration
variable "aws_region" {
  description = "AWS region"
  type        = string
  # No default value to avoid hardcoding
}

# Use data sources to get current region and partition
data "aws_region" "current" {}
data "aws_partition" "current" {}

# Provider configuration without hardcoded region
# Region should come from environment variables (AWS_DEFAULT_REGION) or AWS profiles
provider "aws" {
  # No region specified - will use environment or profile configuration
}

# S3 bucket without hardcoded region
resource "aws_s3_bucket" "example" {
  bucket = "my-example-bucket-${random_id.bucket_suffix.hex}"
}

resource "random_id" "bucket_suffix" {
  byte_length = 4
}

# IAM role with dynamic ARN construction
resource "aws_iam_role" "example" {
  name = "example-role"
  
  assume_role_policy = "{\"Version\": \"2012-10-17\", \"Statement\": [{\"Action\": \"sts:AssumeRole\", \"Effect\": \"Allow\", \"Principal\": {\"Service\": \"ec2.amazonaws.com\"}}]}"
}

# IAM role policy with dynamic ARNs
resource "aws_iam_role_policy" "example" {
  name = "example-policy"
  role = aws_iam_role.example.id

  policy = "{\"Version\": \"2012-10-17\", \"Statement\": [{\"Action\": [\"s3:GetObject\"], \"Effect\": \"Allow\", \"Resource\": \"arn:${data.aws_partition.current.partition}:s3:::my-bucket/${data.aws_region.current.name}/*\"}]}"
}

# IAM policy with dynamic ARNs
resource "aws_iam_policy" "example" {
  name = "example-policy"

  policy = "{\"Version\": \"2012-10-17\", \"Statement\": [{\"Action\": [\"lambda:InvokeFunction\"], \"Effect\": \"Allow\", \"Resource\": \"arn:${data.aws_partition.current.partition}:lambda:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:function:my-function\"}]}"
}

data "aws_caller_identity" "current" {}


# Lambda permission with dynamic ARN (best practice)
resource "aws_lambda_permission" "example" {
  statement_id  = "AllowS3Invoke"
  action        = "lambda:InvokeFunction"
  function_name = "my-function"
  principal     = "s3.amazonaws.com"
  source_arn    = "arn:${data.aws_partition.current.partition}:s3:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:bucket/my-bucket"
}

# SNS subscription with dynamic ARN (best practice)
resource "aws_sns_topic_subscription" "example" {
  topic_arn = "arn:${data.aws_partition.current.partition}:sns:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:my-topic"
  protocol  = "email"
  endpoint  = "example@example.com"
}

# KMS grant with dynamic ARN (best practice)
resource "aws_kms_grant" "example" {
  name              = "my-grant"
  key_id            = "arn:${data.aws_partition.current.partition}:kms:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:key/12345678-1234-1234-1234-123456789012"
  grantee_principal = "arn:${data.aws_partition.current.partition}:iam::${data.aws_caller_identity.current.account_id}:role/my-role"
}

# CloudWatch event target with dynamic ARN (best practice)
resource "aws_cloudwatch_event_target" "example" {
  rule = "my-rule"
  arn  = "arn:${data.aws_partition.current.partition}:lambda:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:function:my-function"
}
