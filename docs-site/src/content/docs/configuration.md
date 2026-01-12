---
title: Configuration
description: How to configure and run the TFLint AWS Meta Ruleset
---

## Running TFLint

After installation and configuration, run TFLint as usual:

```bash
tflint
```

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

## Verifying Configuration

You can verify the plugin is working by running it on the example configurations:

```bash
# Should show 0 issues from our rules
cd examples/passing && tflint

# Should show multiple issues from our rules  
cd examples/failing && tflint
```

## Default Rules

Two rules are enabled by default when you install the plugin:

- **`aws_meta_hardcoded`** - Comprehensive ARN validation across all AWS resources
- **`aws_service_principal_dns_suffix`** - Detects dns_suffix interpolation in service principals

All other rules are disabled by default to avoid overwhelming existing codebases with violations.

## Rule Categories

### Comprehensive Rules (Enabled by Default)
- `aws_meta_hardcoded` - Checks all AWS resources for hardcoded regions/partitions in ARNs
- `aws_service_principal_dns_suffix` - Detects dns_suffix interpolation

### IAM Policy Rules (Disabled by Default)
- `aws_iam_policy_hardcoded_region` - Hardcoded regions in IAM policies
- `aws_iam_policy_hardcoded_partition` - Hardcoded partitions in IAM policies
- `aws_iam_role_policy_hardcoded_region` - Hardcoded regions in IAM role policies
- `aws_iam_role_policy_hardcoded_partition` - Hardcoded partitions in IAM role policies

### Provider Rules (Disabled by Default)
- `aws_provider_hardcoded_region` - Hardcoded regions in provider configuration

### Service Principal Rules (Mixed)
- `aws_service_principal_hardcoded` - Hardcoded DNS suffixes (disabled by default)
- `aws_service_principal_dns_suffix` - DNS suffix interpolation (enabled by default)

## Common Workflows

### Gradual Adoption
Start with default rules and gradually enable more:

1. **Phase 1**: Use default configuration
2. **Phase 2**: Enable IAM rules for new policies
3. **Phase 3**: Enable provider rules
4. **Phase 4**: Enable comprehensive validation

### Legacy Codebase Integration
For existing codebases with many violations:

1. Start with minimal configuration
2. Fix violations incrementally
3. Enable additional rules as violations are resolved
4. Use selective configuration to focus on specific areas

### CI/CD Integration
Add TFLint to your pipeline:

```yaml
# GitHub Actions example
- name: Run TFLint
  run: |
    tflint --init
    tflint
```

Make sure your `.tflint.hcl` is committed to your repository for consistent results across environments.