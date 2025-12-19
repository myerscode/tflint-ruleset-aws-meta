---
title: aws_iam_policy_hardcoded_region
description: Detects hardcoded AWS regions in aws_iam_policy resources.
---

This rule checks `aws_iam_policy` resources for hardcoded AWS regions in policy documents. Similar to the role policy rule, it examines:

- Hardcoded regions in ARNs within policy statements
- Direct region references in policy JSON

## Example violations

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

## Recommended fix

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
