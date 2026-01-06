---
title: IAM Role Policy Hardcoded Regions
description: Detects hardcoded AWS regions in aws_iam_role_policy resources.
ruleName: aws_iam_role_policy_hardcoded_region
---

**Rule:** `aws_iam_role_policy_hardcoded_region`

This rule checks `aws_iam_role_policy` resources for hardcoded AWS regions in policy documents. It examines both JSON policy strings and structured policy documents to detect:

- Hardcoded regions in ARNs within policy statements (e.g., `arn:aws:s3:::bucket/us-east-1/*`)
- Direct region references in policy JSON

## Example violations

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

## Recommended fix

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

## Enabling this rule

This rule is **disabled by default**. To enable it, add it to your `.tflint.hcl`:

```hcl
rule "aws_iam_role_policy_hardcoded_region" {
  enabled = true
}
```
