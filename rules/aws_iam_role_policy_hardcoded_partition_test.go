package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsIamRolePolicyHardcodedPartitionRule(t *testing.T) {
	tests := []struct {
		Name          string
		Content       string
		ExpectedCount int
	}{
		{
			Name: "hardcoded partition in policy",
			Content: `
resource "aws_iam_role_policy" "example" {
  name = "example-policy"
  role = "example-role"
  policy = "arn:aws:s3:::my-bucket/*"
}`,
			ExpectedCount: 1,
		},
		{
			Name: "no hardcoded partitions",
			Content: `
resource "aws_iam_role_policy" "example" {
  name = "example-policy"
  role = "example-role"
  policy = "some-policy-without-arn"
}`,
			ExpectedCount: 0,
		},
	}

	rule := NewAwsIamRolePolicyHardcodedPartitionRule()

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
