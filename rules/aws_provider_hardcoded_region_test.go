package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsProviderHardcodedRegionRule(t *testing.T) {
	tests := []struct {
		Name          string
		Content       string
		ExpectedCount int
	}{
		{
			Name: "hardcoded region in provider",
			Content: `
provider "aws" {
  region = "us-east-1"
}`,
			ExpectedCount: 1,
		},
		{
			Name: "hardcoded region in assume_role ARN",
			Content: `
provider "aws" {
  assume_role {
    role_arn = "arn:aws:iam::123456789012:role/terraform-role"
  }
}`,
			ExpectedCount: 0, // No region in this ARN
		},
		{
			Name: "hardcoded region in assume_role ARN with region",
			Content: `
provider "aws" {
  assume_role {
    role_arn = "arn:aws:iam:us-west-2:123456789012:role/terraform-role"
  }
}`,
			ExpectedCount: 1,
		},

		{
			Name: "using profile - no issues",
			Content: `
provider "aws" {
  profile = "default"
}`,
			ExpectedCount: 0,
		},
		{
			Name: "using environment variables - no issues",
			Content: `
provider "aws" {
  # Region will be picked up from AWS_DEFAULT_REGION environment variable
}`,
			ExpectedCount: 0,
		},
		{
			Name: "multiple providers with mixed configurations",
			Content: `
provider "aws" {
  alias  = "us_east"
  region = "us-east-1"
}

provider "aws" {
  alias   = "eu_west"
  profile = "eu-profile"
}`,
			ExpectedCount: 1, // Only the first provider has hardcoded region
		},
	}

	rule := NewAwsProviderHardcodedRegionRule()

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
