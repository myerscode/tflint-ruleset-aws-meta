# Passing Examples

This directory contains Terraform configurations that follow AWS best practices and should pass all TFLint rules.

## Best Practices Demonstrated

1. **Use data sources for dynamic ARN construction** - Leverage `data.aws_region.current` and `data.aws_partition.current` in IAM policy ARNs
2. **Avoid hardcoded regions in IAM policies** - Use variables or data sources instead of hardcoding region names in policy documents
3. **Avoid hardcoded partitions in IAM policies** - Use `data.aws_partition.current.partition` instead of hardcoding partition values like `aws`, `aws-cn`, or `aws-us-gov`
4. **Dynamic policy construction** - Build IAM policy ARNs dynamically to support multi-region and multi-partition deployments
5. **Secure provider configuration** - Use variables, environment variables, or AWS profiles instead of hardcoding regions or credentials
6. **Dynamic resource configuration** - Avoid hardcoded regions in all AWS resources by using variables and data sources

## Running TFLint

```bash
cd examples/passing
tflint
```

Expected output: No issues found (exit code 0)