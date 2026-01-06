---
title: TFLint AWS Meta Ruleset
description: A TFLint ruleset for AWS best practices, focusing on preventing hardcoded values and promoting flexible, maintainable Terraform code.
---

# TFLint AWS Meta Ruleset

[![Build Status](https://github.com/myerscode/tflint-ruleset-aws-meta/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/myerscode/tflint-ruleset-aws-meta/actions)

A TFLint ruleset for AWS best practices, focusing on preventing hardcoded values and promoting flexible, maintainable Terraform code.

This ruleset helps enforce multi-region and multi-partition compatibility by detecting hardcoded AWS regions and partitions in your Terraform configurations. It provides comprehensive coverage across IAM policies, provider configurations, and all AWS resource types where hardcoded values prevent flexible deployments.

## Why These Rules Matter

**Multi-Region Deployments:** Hardcoded regions prevent your Terraform configurations from being deployed to different AWS regions without modification.

**Multi-Partition Support:** Hardcoded partitions prevent deployment to AWS GovCloud (`aws-us-gov`) or AWS China (`aws-cn`) regions.

**Security:** Hardcoded credentials in provider configurations pose security risks and should be avoided.

**Maintainability:** Dynamic configurations using variables and data sources are easier to maintain and more flexible.

## How It Works

This ruleset uses the [aws-meta](https://github.com/myerscode/aws-meta) Go package to dynamically generate regex patterns for all AWS regions and partitions. Instead of maintaining hardcoded lists, the patterns are built at runtime from the latest AWS metadata.

**Benefits:**
- New AWS regions are automatically detected when the `aws-meta` package is updated
- Covers all AWS partitions (commercial, GovCloud, China, isolated)
- No manual maintenance required for region lists
- Always up-to-date with AWS's latest offerings

## Best Practices

1. **Use data sources:** `data.aws_region.current.name` and `data.aws_partition.current.partition`
2. **Use variables:** Define region and other parameters as variables
3. **Environment variables:** Use `AWS_REGION`, `AWS_PROFILE` environment variables
4. **AWS profiles:** Configure provider to use AWS CLI profiles
5. **IAM roles:** Use IAM roles for authentication instead of hardcoded keys

## Rules Overview

|Name|Description|Severity|Enabled By Default|Link|
| --- | --- | --- | --- | --- |
|aws_meta_hardcoded|Validates that there are no hardcoded AWS regions or partitions in ARN values across all resource types|WARNING|✅|[docs](/rules/aws_meta_hardcoded)|
|aws_iam_role_policy_hardcoded_region|Validates that there are no hardcoded AWS regions in IAM role policy documents|WARNING|❌|[docs](/rules/aws_iam_role_policy_hardcoded_region)|
|aws_iam_role_policy_hardcoded_partition|Validates that there are no hardcoded AWS partitions in IAM role policy documents|WARNING|❌|[docs](/rules/aws_iam_role_policy_hardcoded_partition)|
|aws_iam_policy_hardcoded_region|Validates that there are no hardcoded AWS regions in IAM policy documents|WARNING|❌|[docs](/rules/aws_iam_policy_hardcoded_region)|
|aws_iam_policy_hardcoded_partition|Validates that there are no hardcoded AWS partitions in IAM policy documents|WARNING|❌|[docs](/rules/aws_iam_policy_hardcoded_partition)|
|aws_provider_hardcoded_region|Validates that there are no hardcoded AWS regions in provider configuration|WARNING|❌|[docs](/rules/aws_provider_hardcoded_region)|
|aws_service_principal_hardcoded|Validates that service principals don't use hardcoded DNS suffixes (e.g., amazonaws.com)|WARNING|❌|[docs](/rules/aws_service_principal_hardcoded)|
|aws_service_principal_dns_suffix|Validates that service principals don't use dns_suffix interpolation|WARNING|✅|[docs](/rules/aws_service_principal_dns_suffix)|

For detailed documentation on each rule, see the [Rules](/rules/) section.


