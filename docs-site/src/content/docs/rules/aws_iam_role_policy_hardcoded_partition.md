---
title: IAM Role Policy Hardcoded Partitions
description: Detects hardcoded AWS partitions in aws_iam_role_policy resources.
ruleName: aws_iam_role_policy_hardcoded_partition
---

**Rule:** `aws_iam_role_policy_hardcoded_partition`

This rule checks `aws_iam_role_policy` resources for hardcoded AWS partitions in policy documents. It detects:

- Hardcoded partitions in ARNs (e.g., `arn:aws:`, `arn:aws-cn:`, `arn:aws-us-gov:`)

## Example violations

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

## Recommended fix

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

## Enabling this rule

This rule is **disabled by default**. To enable it, add it to your `.tflint.hcl`:

```hcl
rule "aws_iam_role_policy_hardcoded_partition" {
  enabled = true
}
```
