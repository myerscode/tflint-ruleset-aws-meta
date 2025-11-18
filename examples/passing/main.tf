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

# Provider configuration using variable (best practice)
provider "aws" {
  region = var.aws_region
  # Credentials should come from environment variables, IAM roles, or AWS profiles
  # profile = var.aws_profile  # Alternative: use AWS profile
}

# S3 bucket without hardcoded region
resource "aws_s3_bucket" "example" {
  bucket = "my-example-bucket-${random_id.bucket_suffix.hex}"
}

resource "random_id" "bucket_suffix" {
  byte_length = 4
}

# Instance using data source for availability zone
resource "aws_instance" "example" {
  ami               = "ami-12345678"
  instance_type     = "t2.micro"
  availability_zone = "${data.aws_region.current.name}a"
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

# Additional AWS resources using dynamic region references (best practice)
resource "aws_s3_bucket" "another_example" {
  bucket = "my-another-bucket-${random_id.bucket_suffix2.hex}"
  # No hardcoded region - will use provider region
}

resource "random_id" "bucket_suffix2" {
  byte_length = 4
}

resource "aws_instance" "another_example" {
  ami           = "ami-87654321"
  instance_type = "t2.micro"
  # Using data source for dynamic AZ selection
  availability_zone = "${data.aws_region.current.name}b"
}

resource "aws_db_instance" "another_example" {
  identifier     = "mydb-another"
  engine         = "mysql"
  instance_class = "db.t3.micro"
  # No hardcoded region - will use provider region
}