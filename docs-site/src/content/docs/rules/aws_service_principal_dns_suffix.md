---
title: Service Principal DNS Suffix Interpolation
description: Detects use of `dns_suffix` interpolation in service principals.
ruleName: aws_service_principal_dns_suffix
---

**Rule:** `aws_service_principal_dns_suffix`

This rule checks for use of `dns_suffix` in service principals (e.g., `service.${var.dns_suffix}`) and suggests using `data.aws_service_principal.<name>.name` instead for better maintainability.

It detects both evaluated strings containing `dns_suffix` and raw interpolated expressions that reference `.dns_suffix`.

## Example violations

```hcl
resource "aws_iam_role" "lambda_role" {
  assume_role_policy = jsonencode({
    Statement = [{
      Principal = {
        Service = "lambda.${var.dns_suffix}"  # ❌ Using dns_suffix interpolation
      }
    }]
  })
}

resource "aws_iam_role" "ec2_role" {
  assume_role_policy = jsonencode({
    Statement = [{
      Principal = {
        Service = "ec2.${data.aws_partition.current.dns_suffix}"  # ❌ Using dns_suffix interpolation
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

## Why this matters

The `data.aws_service_principal` data source automatically handles:
- Different DNS suffixes across AWS partitions (`.amazonaws.com`, `.amazonaws.com.cn`, `.amazonaws-us-gov.com`)
- Regional variations in service principals
- Future changes to service principal formats

This approach is more maintainable and partition-aware than manual DNS suffix interpolation.

## Enabling this rule

This rule is **enabled by default** when you install the aws-meta plugin. No additional configuration is needed.

If you want to disable this rule, add it to your `.tflint.hcl`:

```hcl
rule "aws_service_principal_dns_suffix" {
  enabled = false
}
```
