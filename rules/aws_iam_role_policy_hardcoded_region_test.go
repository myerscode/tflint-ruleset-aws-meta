package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsIamRolePolicyHardcodedRegionRule(t *testing.T) {
	tests := []struct {
		Name          string
		Content       string
		ExpectedCount int
	}{
		{
			Name: "hardcoded region in policy",
			Content: `
resource "aws_iam_role_policy" "example" {
  name = "example-policy"
  role = "example-role"
  policy = "arn:aws:s3:::my-bucket/us-east-1/*"
}`,
			ExpectedCount: 1,
		},
		{
			Name: "no hardcoded regions",
			Content: `
resource "aws_iam_role_policy" "example" {
  name = "example-policy"
  role = "example-role"
  policy = "arn:aws:s3:::my-bucket/variable-region/*"
}`,
			ExpectedCount: 0,
		},
	}

	rule := NewAwsIamRolePolicyHardcodedRegionRule()

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
