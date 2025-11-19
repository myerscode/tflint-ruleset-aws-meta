# Failing Examples

This directory contains Terraform configurations that violate AWS best practices and will trigger TFLint rule violations.

## Issues Demonstrated

### aws_iam_role_policy_hardcoded_region violations:
- Hardcoded region in IAM role policy ARNs (`us-east-1`, `us-west-2`, `ap-southeast-1`, `cn-north-1`)

### aws_iam_role_policy_hardcoded_partition violations:
- Hardcoded partition in IAM role policy ARNs (`aws`, `aws-cn`)

### aws_iam_policy_hardcoded_region violations:
- Hardcoded region in IAM policy ARNs (`eu-west-1`, `us-gov-west-1`)

### aws_iam_policy_hardcoded_partition violations:
- Hardcoded partition in IAM policy ARNs (`aws`, `aws-us-gov`)

### aws_provider_hardcoded_region violations:
- Hardcoded region in provider configuration (`us-east-1`, `us-west-2`)
- Hardcoded AWS credentials in provider configuration (access key and secret key)

### aws_hardcoded_region violations:
- Hardcoded region in S3 bucket (`us-west-1`)
- Hardcoded availability zone in EC2 instance (`eu-central-1a`)
- Hardcoded region in RDS instance (`ap-southeast-2`)

## Running TFLint

```bash
cd examples/failing
tflint
```

Expected output: 24 warnings about hardcoded regions, partitions, and credentials across all AWS resources (exit code 2 with warnings)