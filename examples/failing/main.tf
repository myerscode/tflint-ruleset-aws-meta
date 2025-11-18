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

# Hardcoded region in provider (will trigger aws_provider_hardcoded_region rule)
provider "aws" {
  region = "us-east-1"
}

# Provider with hardcoded credentials (will trigger aws_provider_hardcoded_region rule)
provider "aws" {
  alias      = "hardcoded_creds"
  region     = "us-west-2"
  access_key = "AKIAIOSFODNN7EXAMPLE"
  secret_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
}

# S3 bucket with hardcoded region (will trigger aws_hardcoded_region rule)
resource "aws_s3_bucket" "example" {
  bucket = "my-example-bucket"
  region = "us-west-2"
}

# Instance with hardcoded availability zone (will trigger aws_hardcoded_region rule)
resource "aws_instance" "example" {
  ami               = "ami-12345678"
  instance_type     = "t2.micro"
  availability_zone = "eu-west-1a"
}

# IAM role policy with hardcoded region (will trigger aws_iam_role_policy_hardcoded_region rule)
resource "aws_iam_role_policy" "example_region" {
  name = "example-policy-region"
  role = "example-role"

  policy = "{\"Version\": \"2012-10-17\", \"Statement\": [{\"Action\": [\"s3:GetObject\"], \"Effect\": \"Allow\", \"Resource\": \"arn:aws:s3:::my-bucket/us-east-1/*\"}]}"
}

# IAM role policy with hardcoded partition (will trigger aws_iam_role_policy_hardcoded_partition rule)
resource "aws_iam_role_policy" "example_partition" {
  name = "example-policy-partition"
  role = "example-role"

  policy = "{\"Version\": \"2012-10-17\", \"Statement\": [{\"Action\": [\"lambda:InvokeFunction\"], \"Effect\": \"Allow\", \"Resource\": \"arn:aws:lambda:us-west-2:123456789012:function:my-function\"}]}"
}

# IAM policy with hardcoded region (will trigger aws_iam_policy_hardcoded_region rule)
resource "aws_iam_policy" "example_region" {
  name = "example-policy-region"

  policy = "{\"Version\": \"2012-10-17\", \"Statement\": [{\"Action\": [\"dynamodb:GetItem\"], \"Effect\": \"Allow\", \"Resource\": \"arn:aws:dynamodb:eu-west-1:123456789012:table/my-table\"}]}"
}

# IAM policy with hardcoded partition (will trigger aws_iam_policy_hardcoded_partition rule)
resource "aws_iam_policy" "example_partition" {
  name = "example-policy-partition"

  policy = "{\"Version\": \"2012-10-17\", \"Statement\": [{\"Action\": [\"sqs:SendMessage\"], \"Effect\": \"Allow\", \"Resource\": \"arn:aws-us-gov:sqs:us-gov-west-1:123456789012:my-queue\"}]}"
}

# Multiple violations in one policy
resource "aws_iam_role_policy" "multiple_violations" {
  name = "multiple-violations"
  role = "example-role"

  policy = "{\"Version\": \"2012-10-17\", \"Statement\": [{\"Action\": [\"s3:GetObject\"], \"Effect\": \"Allow\", \"Resource\": [\"arn:aws:s3:::bucket1/ap-southeast-1/*\", \"arn:aws-cn:s3:::bucket2/cn-north-1/*\"]}]}"
}

# Additional AWS resources with hardcoded regions (will trigger aws_hardcoded_region rule)
resource "aws_s3_bucket" "hardcoded_region" {
  bucket = "my-hardcoded-bucket"
  region = "us-west-1"
}

resource "aws_instance" "hardcoded_az" {
  ami               = "ami-87654321"
  instance_type     = "t2.micro"
  availability_zone = "eu-central-1a"
}

resource "aws_db_instance" "hardcoded_region" {
  identifier     = "mydb-hardcoded"
  engine         = "mysql"
  instance_class = "db.t3.micro"
  region         = "ap-southeast-2"
}