---
title: aws_meta_hardcoded
description: Checks all AWS resources for hardcoded regions and partitions in ARN values.
---

This is a comprehensive rule that checks ALL AWS resources for hardcoded regions and partitions in ARN values. It works by walking through all expressions in your Terraform files and detecting any string that looks like an ARN with hardcoded values.

This rule covers resource types including:

- Lambda (permissions, event source mappings, functions)
- SNS/SQS (subscriptions, queue policies)
- CloudWatch (event targets, log subscriptions, alarms)
- API Gateway (integrations, authorizers)
- KMS (grants, aliases, keys)
- Secrets Manager (rotations, policies)
- ECS (services, task definitions)
- RDS (instances, event subscriptions, clusters)
- S3 (notifications, policies, access points)
- And many more...

## Example violations

```hcl
resource "aws_lambda_permission" "test" {
  source_arn = "arn:aws:s3:us-east-1:123456789012:bucket/my-bucket"  # ❌ Hardcoded region and partition
}

resource "aws_kms_grant" "test" {
  key_id = "arn:aws:kms:eu-west-1:123456789012:key/12345678-1234-1234-1234-123456789012"  # ❌ Hardcoded region and partition
}
```

## Recommended fixes

```hcl
data "aws_region" "current" {}
data "aws_partition" "current" {}

resource "aws_lambda_permission" "test" {
  source_arn = "arn:${data.aws_partition.current.partition}:s3:${data.aws_region.current.name}:123456789012:bucket/my-bucket"  # ✅ Dynamic
}

resource "aws_kms_grant" "test" {
  key_id = "arn:${data.aws_partition.current.partition}:kms:${data.aws_region.current.name}:123456789012:key/12345678-1234-1234-1234-123456789012"  # ✅ Dynamic
}
```
