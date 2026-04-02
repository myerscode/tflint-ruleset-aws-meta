package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsHardcodedIDsRule(t *testing.T) {
	tests := []struct {
		Name          string
		Content       string
		ExpectedCount int
	}{
		{
			Name: "hardcoded account ID in ARN",
			Content: `
resource "aws_iam_role_policy" "test" {
  policy = jsonencode({
    Statement = [{
      Effect   = "Allow"
      Action   = "s3:GetObject"
      Resource = "arn:aws:s3:::my-bucket/123456789012/*"
    }]
  })
}`,
			ExpectedCount: 2,
		},
		{
			Name: "hardcoded account ID as standalone value",
			Content: `
resource "aws_guardduty_member" "test" {
  account_id = "123456789012"
}`,
			ExpectedCount: 2,
		},
		{
			Name: "hardcoded AMI ID",
			Content: `
resource "aws_instance" "test" {
  ami           = "ami-0abcdef1234567890"
  instance_type = "t3.micro"
}`,
			ExpectedCount: 2,
		},
		{
			Name: "hardcoded AMI ID in launch template",
			Content: `
resource "aws_launch_template" "test" {
  image_id = "ami-0ff8a91507f77f867"
}`,
			ExpectedCount: 2,
		},
		{
			Name: "both account ID and AMI in same config",
			Content: `
resource "aws_instance" "test" {
  ami = "ami-0abcdef1234567890"
}

resource "aws_guardduty_member" "test" {
  account_id = "123456789012"
}`,
			ExpectedCount: 4,
		},
		{
			Name: "dynamic account ID using data source",
			Content: `
data "aws_caller_identity" "current" {}

resource "aws_guardduty_member" "test" {
  account_id = data.aws_caller_identity.current.account_id
}`,
			ExpectedCount: 0,
		},
		{
			Name: "dynamic AMI using data source",
			Content: `
data "aws_ami" "ubuntu" {
  most_recent = true
  owners      = ["099720109477"]
}

resource "aws_instance" "test" {
  ami           = data.aws_ami.ubuntu.id
  instance_type = "t3.micro"
}`,
			ExpectedCount: 2, // the owners value is a 12-digit account ID
		},
		{
			Name: "no hardcoded IDs",
			Content: `
resource "aws_s3_bucket" "test" {
  bucket = "my-bucket-name"
}`,
			ExpectedCount: 0,
		},
		{
			Name: "short number not an account ID",
			Content: `
resource "aws_instance" "test" {
  tags = {
    Port = "8080"
  }
}`,
			ExpectedCount: 0,
		},
	}

	rule := NewAwsHardcodedIDsRule()

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
