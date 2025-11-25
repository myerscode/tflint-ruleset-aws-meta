package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsServicePrincipalDNSSuffixRule(t *testing.T) {
	tests := []struct {
		Name          string
		Content       string
		ExpectedCount int
	}{
		{
			Name: "using dns_suffix with single service",
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
			ExpectedCount: 7,
		},
		{
			Name: "using dns_suffix with multiple services",
			Content: `
data "aws_partition" "current" {}

resource "aws_iam_role" "test" {
  assume_role_policy = jsonencode({
    Statement = [{
      Principal = {
        Service = [
          "lambda.${data.aws_partition.current.dns_suffix}",
          "ec2.${data.aws_partition.current.dns_suffix}",
          "ecs-tasks.${data.aws_partition.current.dns_suffix}"
        ]
      }
    }]
  })
}`,
			ExpectedCount: 12,
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
			Name: "hardcoded service principal (no issues for this rule)",
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
			ExpectedCount: 0,
		},
	}

	rule := NewAwsServicePrincipalDNSSuffixRule()

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
