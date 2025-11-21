package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsMetaHardcodedRule(t *testing.T) {
	tests := []struct {
		Name          string
		Content       string
		ExpectedCount int
	}{
		{
			Name: "iam role assume_role_policy with hardcoded ARN",
			Content: `
resource "aws_iam_role" "test" {
  assume_role_policy = jsonencode({
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        AWS = "arn:aws:iam:us-east-1:123456789012:root"
      }
    }]
  })
}`,
			ExpectedCount: 4, // 2x because WalkExpressions visits nested expressions
		},
		{
			Name: "lambda permission with hardcoded source_arn",
			Content: `
resource "aws_lambda_permission" "test" {
  source_arn = "arn:aws:s3:eu-west-1:123456789012:bucket/my-bucket"
}`,
			ExpectedCount: 4,
		},
		{
			Name: "lambda event source mapping with hardcoded event_source_arn",
			Content: `
resource "aws_lambda_event_source_mapping" "test" {
  event_source_arn = "arn:aws:dynamodb:us-east-1:123456789012:table/my-table"
}`,
			ExpectedCount: 4,
		},
		{
			Name: "sns subscription with hardcoded topic_arn",
			Content: `
resource "aws_sns_topic_subscription" "test" {
  topic_arn = "arn:aws:sns:us-west-2:123456789012:my-topic"
}`,
			ExpectedCount: 4,
		},
		{
			Name: "cloudwatch event target with hardcoded arn",
			Content: `
resource "aws_cloudwatch_event_target" "test" {
  arn = "arn:aws:lambda:us-west-2:123456789012:function:my-function"
}`,
			ExpectedCount: 4,
		},
		{
			Name: "cloudwatch log subscription filter with hardcoded destination_arn",
			Content: `
resource "aws_cloudwatch_log_subscription_filter" "test" {
  destination_arn = "arn:aws:lambda:eu-west-1:123456789012:function:my-function"
}`,
			ExpectedCount: 4,
		},
		{
			Name: "api gateway integration with hardcoded uri",
			Content: `
resource "aws_api_gateway_integration" "test" {
  uri = "arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/arn:aws:lambda:us-east-1:123456789012:function:my-function/invocations"
}`,
			ExpectedCount: 4, // Only the lambda ARN is detected, not the apigateway ARN format
		},
		{
			Name: "kms grant with hardcoded key_id",
			Content: `
resource "aws_kms_grant" "test" {
  key_id = "arn:aws:kms:us-west-2:123456789012:key/12345678-1234-1234-1234-123456789012"
}`,
			ExpectedCount: 4,
		},
		{
			Name: "kms alias with hardcoded target_key_id",
			Content: `
resource "aws_kms_alias" "test" {
  target_key_id = "arn:aws:kms:eu-west-1:123456789012:key/12345678-1234-1234-1234-123456789012"
}`,
			ExpectedCount: 4,
		},
		{
			Name: "secretsmanager rotation with hardcoded rotation_lambda_arn",
			Content: `
resource "aws_secretsmanager_secret_rotation" "test" {
  rotation_lambda_arn = "arn:aws:lambda:us-east-1:123456789012:function:my-rotation-function"
}`,
			ExpectedCount: 4,
		},
		{
			Name: "db instance with hardcoded replicate_source_db",
			Content: `
resource "aws_db_instance" "test" {
  replicate_source_db = "arn:aws:rds:us-east-1:123456789012:db:my-source-db"
}`,
			ExpectedCount: 4,
		},
		{
			Name: "db event subscription with hardcoded sns_topic",
			Content: `
resource "aws_db_event_subscription" "test" {
  sns_topic = "arn:aws:sns:eu-west-1:123456789012:my-topic"
}`,
			ExpectedCount: 4,
		},
		{
			Name: "multiple resources with different partitions",
			Content: `
resource "aws_lambda_permission" "test1" {
  source_arn = "arn:aws:s3:us-east-1:123456789012:bucket/my-bucket"
}

resource "aws_sns_topic_subscription" "test2" {
  topic_arn = "arn:aws-cn:sns:cn-north-1:123456789012:my-topic"
}`,
			ExpectedCount: 8,
		},
		{
			Name: "resource with dynamic ARN using data sources",
			Content: `
data "aws_region" "current" {}
data "aws_partition" "current" {}

resource "aws_lambda_permission" "test" {
  source_arn = "arn:${data.aws_partition.current.partition}:s3:${data.aws_region.current.name}:123456789012:bucket/my-bucket"
}`,
			ExpectedCount: 0,
		},

		{
			Name: "non-ARN string values",
			Content: `
resource "aws_s3_bucket" "test" {
  bucket = "my-bucket-name"
}`,
			ExpectedCount: 0,
		},
	}

	rule := NewAwsMetaHardcodedRule()

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"main.tf": test.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			if len(runner.Issues) != test.ExpectedCount {
				t.Errorf("Expected %d issues, got %d", test.ExpectedCount, len(runner.Issues))
				for i, issue := range runner.Issues {
					t.Logf("Issue %d: %s", i+1, issue.Message)
				}
			}
		})
	}
}
