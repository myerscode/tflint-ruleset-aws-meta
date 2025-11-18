# TFLint AWS Meta Ruleset

[![Build Status](https://github.com/myerscode/tflint-ruleset-aws-meta/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/myerscode/tflint-ruleset-aws-meta/actions)

A TFLint ruleset for AWS best practices, focusing on preventing hardcoded values and promoting flexible, maintainable Terraform code.

This ruleset helps enforce multi-region and multi-partition compatibility by detecting hardcoded AWS regions and partitions in your Terraform configurations. It provides comprehensive coverage across IAM policies, provider configurations, and all AWS resource types where hardcoded values prevent flexible deployments.

## Requirements

- TFLint v0.42+
- Go v1.25

## Installation

TODO: This template repository does not contain release binaries, so this installation will not work. Please rewrite for your repository. See the "Building the plugin" section to get this template ruleset working.

You can install the plugin with `tflint --init`. Declare a config in `.tflint.hcl` as follows:

```hcl
plugin "aws-multi" {
  enabled = true

  version = "0.1.0"
  source  = "github.com/myerscode/tflint-ruleset-aws-meta"
}
```

## Rules

|Name|Description|Severity|Enabled|Link|
| --- | --- | --- | --- | --- |
|aws_iam_role_policy_hardcoded_region|Validates that there are no hardcoded AWS regions in IAM role policy documents|WARNING|✔||
|aws_iam_role_policy_hardcoded_partition|Validates that there are no hardcoded AWS partitions in IAM role policy documents|WARNING|✔||
|aws_iam_policy_hardcoded_region|Validates that there are no hardcoded AWS regions in IAM policy documents|WARNING|✔||
|aws_iam_policy_hardcoded_partition|Validates that there are no hardcoded AWS partitions in IAM policy documents|WARNING|✔||
|aws_provider_hardcoded_region|Validates that there are no hardcoded AWS regions or credentials in provider configuration|WARNING|✔||

### Rule Details

#### aws_iam_role_policy_hardcoded_region

This rule checks `aws_iam_role_policy` resources for hardcoded AWS regions in policy documents. It examines both JSON policy strings and structured policy documents to detect:

- Hardcoded regions in ARNs within policy statements (e.g., `arn:aws:s3:::bucket/us-east-1/*`)
- Direct region references in policy JSON

**Example violations:**
```hcl
resource "aws_iam_role_policy" "bad" {
  policy = jsonencode({
    Statement = [{
      Effect   = "Allow"
      Action   = "s3:GetObject"
      Resource = "arn:aws:s3:::bucket/us-east-1/*"  # ❌ Hardcoded region
    }]
  })
}
```

**Recommended fix:**
```hcl
resource "aws_iam_role_policy" "good" {
  policy = jsonencode({
    Statement = [{
      Effect   = "Allow"
      Action   = "s3:GetObject"
      Resource = "arn:aws:s3:::bucket/${data.aws_region.current.name}/*"  # ✅ Dynamic region
    }]
  })
}
```

#### aws_iam_role_policy_hardcoded_partition

This rule checks `aws_iam_role_policy` resources for hardcoded AWS partitions in policy documents. It detects:

- Hardcoded partitions in ARNs (e.g., `arn:aws:`, `arn:aws-cn:`, `arn:aws-us-gov:`)

**Example violations:**
```hcl
resource "aws_iam_role_policy" "bad" {
  policy = jsonencode({
    Statement = [{
      Effect   = "Allow"
      Action   = "s3:*"
      Resource = "arn:aws:s3:::bucket/*"  # ❌ Hardcoded partition
    }]
  })
}
```

**Recommended fix:**
```hcl
resource "aws_iam_role_policy" "good" {
  policy = jsonencode({
    Statement = [{
      Effect   = "Allow"
      Action   = "s3:*"
      Resource = "arn:${data.aws_partition.current.partition}:s3:::bucket/*"  # ✅ Dynamic partition
    }]
  })
}
```

#### aws_iam_policy_hardcoded_region

This rule checks `aws_iam_policy` resources for hardcoded AWS regions in policy documents. Similar to the role policy rule, it examines:

- Hardcoded regions in ARNs within policy statements
- Direct region references in policy JSON

**Example violations:**
```hcl
resource "aws_iam_policy" "bad" {
  policy = jsonencode({
    Statement = [{
      Effect   = "Allow"
      Action   = "lambda:InvokeFunction"
      Resource = "arn:aws:lambda:eu-west-1:123456789012:function:*"  # ❌ Hardcoded region
    }]
  })
}
```

**Recommended fix:**
```hcl
resource "aws_iam_policy" "good" {
  policy = jsonencode({
    Statement = [{
      Effect   = "Allow"
      Action   = "lambda:InvokeFunction"
      Resource = "arn:aws:lambda:${data.aws_region.current.name}:123456789012:function:*"  # ✅ Dynamic region
    }]
  })
}
```

#### aws_iam_policy_hardcoded_partition

This rule checks `aws_iam_policy` resources for hardcoded AWS partitions in policy documents. It detects:

- Hardcoded partitions in ARNs within policy statements

**Example violations:**
```hcl
resource "aws_iam_policy" "bad" {
  policy = jsonencode({
    Statement = [{
      Effect   = "Allow"
      Action   = "sqs:*"
      Resource = "arn:aws-us-gov:sqs:*:*:*"  # ❌ Hardcoded partition
    }]
  })
}
```

**Recommended fix:**
```hcl
resource "aws_iam_policy" "good" {
  policy = jsonencode({
    Statement = [{
      Effect   = "Allow"
      Action   = "sqs:*"
      Resource = "arn:${data.aws_partition.current.partition}:sqs:*:*:*"  # ✅ Dynamic partition
    }]
  })
}
```

#### aws_provider_hardcoded_region

This rule checks AWS provider configurations for security and flexibility issues. It detects:

- Hardcoded regions in provider `region` attribute
- Hardcoded AWS access keys and secret keys (security risk)
- Hardcoded regions in `assume_role` ARNs

**Example violations:**
```hcl
provider "aws" {
  region = "us-east-1"  # ❌ Hardcoded region
}
```

**Recommended fix:**
```hcl
provider "aws" {
  region = var.aws_region  # ✅ Use variables
}
```


