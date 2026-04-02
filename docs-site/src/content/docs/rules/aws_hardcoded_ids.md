---
title: Hardcoded AWS IDs
description: Detects hardcoded AWS account IDs and AMI IDs.
ruleName: aws_hardcoded_ids
---

**Rule:** `aws_hardcoded_ids`

This rule checks for hardcoded AWS account IDs and AMI IDs across all expressions in Terraform files.

- **Account IDs** are 12-digit numbers that should be dynamically resolved using `data.aws_caller_identity.current.account_id` or passed as variables.
- **AMI IDs** are region-specific and should be dynamically resolved using `data.aws_ami` lookups.

## Example violations

```hcl
resource "aws_instance" "web" {
  ami           = "ami-0abcdef1234567890"  # ❌ Hardcoded AMI ID
  instance_type = "t3.micro"
}

resource "aws_guardduty_member" "member" {
  account_id = "123456789012"  # ❌ Hardcoded account ID
}

resource "aws_iam_role_policy" "example" {
  policy = jsonencode({
    Statement = [{
      Effect   = "Allow"
      Action   = "sts:AssumeRole"
      Resource = "arn:aws:iam::123456789012:role/my-role"  # ❌ Hardcoded account ID in ARN
    }]
  })
}
```

## Recommended fixes

```hcl
data "aws_caller_identity" "current" {}

data "aws_ami" "ubuntu" {
  most_recent = true
  owners      = ["099720109477"]

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-*-amd64-server-*"]
  }
}

resource "aws_instance" "web" {
  ami           = data.aws_ami.ubuntu.id  # ✅ Dynamic AMI lookup
  instance_type = "t3.micro"
}

resource "aws_guardduty_member" "member" {
  account_id = data.aws_caller_identity.current.account_id  # ✅ Dynamic account ID
}

resource "aws_iam_role_policy" "example" {
  policy = jsonencode({
    Statement = [{
      Effect   = "Allow"
      Action   = "sts:AssumeRole"
      Resource = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/my-role"  # ✅ Dynamic
    }]
  })
}
```

## Enabling this rule

This rule is **disabled by default**. To enable it, add the following to your `.tflint.hcl`:

```hcl
rule "aws_hardcoded_ids" {
  enabled = true
}
```
