# Examples

This directory contains example Terraform configurations to demonstrate the TFLint AWS best practices ruleset.

## Structure

- `passing/` - Examples that follow best practices and should pass all rules
- `failing/` - Examples that violate best practices and will trigger rule violations

## Usage

To test the rules against these examples:

1. Build the plugin:
   ```bash
   make
   ```

2. Install the plugin:
   ```bash
   make install
   ```

3. Test against passing examples:
   ```bash
   cd examples/passing
   tflint
   ```
   This should produce no warnings or errors.

4. Test against failing examples:
   ```bash
   cd examples/failing
   tflint
   ```
   This should produce multiple warnings about hardcoded regions and partitions.

## Rules Demonstrated

### aws_iam_role_policy_hardcoded_region
- Detects hardcoded AWS regions in IAM role policy documents
- Checks both JSON policy strings and ARN patterns within policies

### aws_iam_role_policy_hardcoded_partition
- Detects hardcoded AWS partitions in IAM role policy documents
- Identifies hardcoded partitions in ARNs within policy statements

### aws_iam_policy_hardcoded_region
- Detects hardcoded AWS regions in IAM policy documents
- Checks both JSON policy strings and ARN patterns within policies

### aws_iam_policy_hardcoded_partition
- Detects hardcoded AWS partitions in IAM policy documents
- Identifies hardcoded partitions in ARNs within policy statements

### aws_provider_hardcoded_region
- Detects hardcoded AWS regions in provider configuration
- Detects hardcoded credentials (access keys and secret keys) in provider configuration
- Checks assume_role ARNs for hardcoded regions

### aws_hardcoded_region
- Detects hardcoded AWS regions across all AWS resource types
- Checks region attributes, availability zones, and region-specific identifiers
- Covers S3 buckets, EC2 instances, RDS instances, Lambda functions, and many more AWS services