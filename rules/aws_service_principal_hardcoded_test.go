package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsServicePrincipalHardcodedRule(t *testing.T) {
	tests := []struct {
		Name          string
		Content       string
		ExpectedCount int
	}{
		{
			Name: "hardcoded service principal amazonaws.com",
			Content: `
resource "aws_iam_role" "test" {
  assume_role_policy = jsonencode({
    Statement = [{
      Principal = {
        Service = "s3.amazonaws.com"
      }
    }]
  })
}`,
			ExpectedCount: 2,
		},
		{
			Name: "hardcoded service principal amazonaws.com.cn",
			Content: `
resource "aws_iam_role" "test" {
  assume_role_policy = jsonencode({
    Statement = [{
      Principal = {
        Service = "lambda.amazonaws.com.cn"
      }
    }]
  })
}`,
			ExpectedCount: 2,
		},
		{
			Name: "hardcoded service principal amazonaws-us-gov.com",
			Content: `
resource "aws_iam_role" "test" {
  assume_role_policy = jsonencode({
    Statement = [{
      Principal = {
        Service = "ec2.amazonaws-us-gov.com"
      }
    }]
  })
}`,
			ExpectedCount: 2,
		},
		{
			Name: "multiple hardcoded service principals",
			Content: `
resource "aws_iam_role" "test" {
  assume_role_policy = jsonencode({
    Statement = [{
      Principal = {
        Service = [
          "lambda.amazonaws.com",
          "ec2.amazonaws.com",
          "ecs-tasks.amazonaws.com"
        ]
      }
    }]
  })
}`,
			ExpectedCount: 6,
		},
		{
			Name: "using data source (no issues)",
			Content: `
data "aws_service_principal" "s3" {
  service_name = "s3"
}

resource "aws_iam_role" "test" {
  assume_role_policy = jsonencode({
    Statement = [{
      Principal = {
        Service = data.aws_service_principal.s3.name
      }
    }]
  })
}`,
			ExpectedCount: 0,
		},
		{
			Name: "using dns_suffix (no issues for this rule)",
			Content: `
data "aws_partition" "current" {}

resource "aws_iam_role" "test" {
  assume_role_policy = jsonencode({
    Statement = [{
      Principal = {
        Service = "s3.${data.aws_partition.current.dns_suffix}"
      }
    }]
  })
}`,
			ExpectedCount: 0,
		},
	}

	rule := NewAwsServicePrincipalHardcodedRule()

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
