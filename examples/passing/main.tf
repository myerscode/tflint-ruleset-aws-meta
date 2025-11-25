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

# Get service principal for EC2
data "aws_service_principal" "ec2_for_role" {
  service_name = "ec2"
}

# IAM role with dynamic service principal (best practice)
resource "aws_iam_role" "example" {
  name = "example-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = data.aws_service_principal.ec2_for_role.name
      }
    }]
  })
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


# Get service principal for S3 (for Lambda permission)
data "aws_service_principal" "s3_for_lambda" {
  service_name = "s3"
}

# Lambda permission with dynamic ARN and service principal (best practice)
resource "aws_lambda_permission" "example" {
  statement_id  = "AllowS3Invoke"
  action        = "lambda:InvokeFunction"
  function_name = "my-function"
  principal     = data.aws_service_principal.s3_for_lambda.name
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
  operations        = ["Encrypt", "Decrypt", "GenerateDataKey"]
  key_id            = "arn:${data.aws_partition.current.partition}:kms:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:key/12345678-1234-1234-1234-123456789012"
  grantee_principal = "arn:${data.aws_partition.current.partition}:iam::${data.aws_caller_identity.current.account_id}:role/my-role"
}

# CloudWatch event target with dynamic ARN (best practice)
resource "aws_cloudwatch_event_target" "example" {
  rule = "my-rule"
  arn  = "arn:${data.aws_partition.current.partition}:lambda:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:function:my-function"
}

# Get service principal names dynamically (best practice)
data "aws_service_principal" "s3" {
  service_name = "s3"
}

data "aws_service_principal" "lambda" {
  service_name = "lambda"
}

data "aws_service_principal" "ec2" {
  service_name = "ec2"
}

data "aws_service_principal" "ecs_tasks" {
  service_name = "ecs-tasks"
}

# IAM role with service principal using data source (best practice)
resource "aws_iam_role" "service_principal_example" {
  name = "service-principal-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = data.aws_service_principal.s3.name
      }
    }]
  })
}

# IAM role with multiple service principals using data sources (best practice)
resource "aws_iam_role" "multi_service_principal" {
  name = "multi-service-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = [
          data.aws_service_principal.lambda.name,
          data.aws_service_principal.ec2.name,
          data.aws_service_principal.ecs_tasks.name
        ]
      }
    }]
  })
}
