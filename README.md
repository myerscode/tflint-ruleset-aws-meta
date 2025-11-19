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
plugin "aws-meta" {
  enabled = true

  version = "0.1.0"
  source  = "github.com/myerscode/tflint-ruleset-aws-meta"
}
```

## Rules

|Name|Description|Severity|Enabled|Link|
| --- | --- | --- | --- | --- |
|aws_iam_role_policy_hardcoded_region|Validates that there are no hardcoded AWS regions in IAM role policy documents|WARNING|✔|[docs](docs/rules.md#aws_iam_role_policy_hardcoded_region)|
|aws_iam_role_policy_hardcoded_partition|Validates that there are no hardcoded AWS partitions in IAM role policy documents|WARNING|✔|[docs](docs/rules.md#aws_iam_role_policy_hardcoded_partition)|
|aws_iam_policy_hardcoded_region|Validates that there are no hardcoded AWS regions in IAM policy documents|WARNING|✔|[docs](docs/rules.md#aws_iam_policy_hardcoded_region)|
|aws_iam_policy_hardcoded_partition|Validates that there are no hardcoded AWS partitions in IAM policy documents|WARNING|✔|[docs](docs/rules.md#aws_iam_policy_hardcoded_partition)|
|aws_provider_hardcoded_region|Validates that there are no hardcoded AWS regions in provider configuration|WARNING|✔|[docs](docs/rules.md#aws_provider_hardcoded_region)|
|aws_arn_hardcoded|Validates that there are no hardcoded AWS regions or partitions in ARN values across all resource types|WARNING|✔|[docs](docs/rules.md#aws_arn_hardcoded)|

For detailed examples and usage information, see the [Rule Details documentation](docs/rules.md).


