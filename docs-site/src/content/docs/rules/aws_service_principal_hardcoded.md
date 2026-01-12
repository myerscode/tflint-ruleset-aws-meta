---
title: Hardcoded Service Principal DNS Suffixes
description: Detects hardcoded AWS service principal DNS suffixes.
ruleName: aws_service_principal_hardcoded
---

**Rule:** `aws_service_principal_hardcoded`

This rule checks for hardcoded AWS service principal DNS suffixes in expressions and strings across Terraform files.

It matches service principal forms like `service.amazonaws.com`, `service.amazonaws.com.cn`, and `service.amazonaws-us-gov.com` and emits an issue suggesting use of data sources (e.g., `data.aws_service_principal.<name>.name`) for multi-partition compatibility.

**Note:** This rule is disabled by default as it may produce many findings in existing codebases. Enable it when you want to migrate to using data sources.

## Example violations

```hcl
resource "aws_iam_role" "lambda_role" {
  assume_role_policy = jsonencode({
    Statement = [{
      Principal = {
        Service = "lambda.amazonaws.com"  # ❌ Hardcoded DNS suffix
      }
    }]
  })
}

resource "aws_iam_role" "ec2_role" {
  assume_role_policy = jsonencode({
    Statement = [{
      Principal = {
        Service = "ec2.amazonaws-us-gov.com"  # ❌ Hardcoded GovCloud DNS suffix
      }
    }]
  })
}
```

## Recommended fixes

```hcl
data "aws_service_principal" "lambda" {
  service_name = "lambda"
}

data "aws_service_principal" "ec2" {
  service_name = "ec2"
}

resource "aws_iam_role" "lambda_role" {
  assume_role_policy = jsonencode({
    Statement = [{
      Principal = {
        Service = data.aws_service_principal.lambda.name  # ✅ Using data source
      }
    }]
  })
}

resource "aws_iam_role" "ec2_role" {
  assume_role_policy = jsonencode({
    Statement = [{
      Principal = {
        Service = data.aws_service_principal.ec2.name  # ✅ Using data source
      }
    }]
  })
}
```

## Enabling this rule

To enable this optional rule, add it to your `.tflint.hcl`:

```hcl
rule "aws_service_principal_hardcoded" {
  enabled = true
}
```
