---
title: Installation
description: How to install and configure the TFLint AWS Meta Ruleset
---

## Requirements

- TFLint v0.42+
- Go v1.25

## Installing the Plugin

TODO: This template repository does not contain release binaries, so this installation will not work. Please rewrite for your repository. See the "Building the plugin" section to get this template ruleset working.

You can install the plugin with `tflint --init`. Declare a config in `.tflint.hcl` as follows:

```hcl
plugin "aws-meta" {
  enabled = true

  version = "0.1.0"
  source  = "github.com/myerscode/tflint-ruleset-aws-meta"
}
```

## Basic Configuration

Once installed, the plugin will run with default settings. Two rules are enabled by default:

- `aws_meta_hardcoded` - Comprehensive ARN validation across all AWS resources
- `aws_service_principal_dns_suffix` - Detects dns_suffix interpolation in service principals

## Enabling Additional Rules

Most rules are disabled by default to avoid overwhelming existing codebases. You can enable specific rules by adding them to your `.tflint.hcl`:

```hcl
plugin "aws-meta" {
  enabled = true
  version = "0.1.0"
  source  = "github.com/myerscode/tflint-ruleset-aws-meta"
}

# Enable specific rules
rule "aws_iam_policy_hardcoded_region" {
  enabled = true
}

rule "aws_iam_policy_hardcoded_partition" {
  enabled = true
}

rule "aws_provider_hardcoded_region" {
  enabled = true
}
```

## Disabling Rules

You can disable any rule, including the default ones:

```hcl
# Disable a default rule
rule "aws_meta_hardcoded" {
  enabled = false
}
```

## Running TFLint

After configuration, run TFLint as usual:

```bash
tflint
```

## Verifying Installation

You can verify the plugin is working by running it on the example configurations:

```bash
# Should show 0 issues from our rules
cd examples/passing && tflint

# Should show multiple issues from our rules  
cd examples/failing && tflint
```

## Configuration Examples

### Minimal Configuration (Default Rules Only)
```hcl
plugin "aws-meta" {
  enabled = true
  version = "0.1.0"
  source  = "github.com/myerscode/tflint-ruleset-aws-meta"
}
```

### Comprehensive Configuration (All Rules Enabled)
```hcl
plugin "aws-meta" {
  enabled = true
  version = "0.1.0"
  source  = "github.com/myerscode/tflint-ruleset-aws-meta"
}

rule "aws_iam_role_policy_hardcoded_region" {
  enabled = true
}

rule "aws_iam_role_policy_hardcoded_partition" {
  enabled = true
}

rule "aws_iam_policy_hardcoded_region" {
  enabled = true
}

rule "aws_iam_policy_hardcoded_partition" {
  enabled = true
}

rule "aws_provider_hardcoded_region" {
  enabled = true
}

rule "aws_service_principal_hardcoded" {
  enabled = true
}
```

### Selective Configuration (IAM Rules Only)
```hcl
plugin "aws-meta" {
  enabled = true
  version = "0.1.0"
  source  = "github.com/myerscode/tflint-ruleset-aws-meta"
}

# Disable the comprehensive rule
rule "aws_meta_hardcoded" {
  enabled = false
}

# Enable specific IAM rules
rule "aws_iam_policy_hardcoded_region" {
  enabled = true
}

rule "aws_iam_policy_hardcoded_partition" {
  enabled = true
}

rule "aws_iam_role_policy_hardcoded_region" {
  enabled = true
}

rule "aws_iam_role_policy_hardcoded_partition" {
  enabled = true
}
```