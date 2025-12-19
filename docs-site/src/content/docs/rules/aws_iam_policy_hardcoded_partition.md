---
title: aws_iam_policy_hardcoded_partition
description: Detects hardcoded AWS partitions in aws_iam_policy resources.
---

This rule checks `aws_iam_policy` resources for hardcoded AWS partitions in policy documents. It detects:

- Hardcoded partitions in ARNs within policy statements

## Example violations

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

## Recommended fix

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
