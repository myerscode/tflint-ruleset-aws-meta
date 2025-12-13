# Rule Details

## aws_meta_hardcoded

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

**Example violations:**
```hcl
resource "aws_lambda_permission" "test" {
  source_arn = "arn:aws:s3:us-east-1:123456789012:bucket/my-bucket"  # ❌ Hardcoded region and partition
}

resource "aws_kms_grant" "test" {
  key_id = "arn:aws:kms:eu-west-1:123456789012:key/12345678-1234-1234-1234-123456789012"  # ❌ Hardcoded region and partition
}
```

**Recommended fixes:**
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

## aws_iam_role_policy_hardcoded_region

This rule checks `aws_iam_role_policy` resources for hardcoded AWS regions in policy documents. It examines both JSON policy strings and structured policy documents to detect:

- Hardcoded regions in ARNs within policy statements (e.g., `arn:aws:s3:::bucket/us-east-1/*`)
- Direct region references in policy JSON

**Example violations:**
```hcl
resource "aws_iam_role_policy" "bad" {
  policy = jsonencode({
    Statement = [{
      Effect   = "Allow"
      Action   = "s3:GetObject"
      Resource = "arn:aws:s3:::bucket/us-east-1/*"  # ❌ Hardcoded region
    }]
  })
}
```

**Recommended fix:**
```hcl
resource "aws_iam_role_policy" "good" {
  policy = jsonencode({
    Statement = [{
      Effect   = "Allow"
      Action   = "s3:GetObject"
      Resource = "arn:aws:s3:::bucket/${data.aws_region.current.name}/*"  # ✅ Dynamic region
    }]
  })
}
```

## aws_iam_role_policy_hardcoded_partition

This rule checks `aws_iam_role_policy` resources for hardcoded AWS partitions in policy documents. It detects:

- Hardcoded partitions in ARNs (e.g., `arn:aws:`, `arn:aws-cn:`, `arn:aws-us-gov:`)

**Example violations:**
```hcl
resource "aws_iam_role_policy" "bad" {
  policy = jsonencode({
    Statement = [{
      Effect   = "Allow"
      Action   = "s3:*"
      Resource = "arn:aws:s3:::bucket/*"  # ❌ Hardcoded partition
    }]
  })
}
```

**Recommended fix:**
```hcl
resource "aws_iam_role_policy" "good" {
  policy = jsonencode({
    Statement = [{
      Effect   = "Allow"
      Action   = "s3:*"
      Resource = "arn:${data.aws_partition.current.partition}:s3:::bucket/*"  # ✅ Dynamic partition
    }]
  })
}
```

## aws_iam_policy_hardcoded_region

This rule checks `aws_iam_policy` resources for hardcoded AWS regions in policy documents. Similar to the role policy rule, it examines:

- Hardcoded regions in ARNs within policy statements
- Direct region references in policy JSON

**Example violations:**
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

**Recommended fix:**
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

## aws_iam_policy_hardcoded_partition

This rule checks `aws_iam_policy` resources for hardcoded AWS partitions in policy documents. It detects:

- Hardcoded partitions in ARNs within policy statements

**Example violations:**
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

**Recommended fix:**
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

## aws_provider_hardcoded_region

This rule checks AWS provider configurations for hardcoded regions. It detects:

- Hardcoded regions in provider `region` attribute
- Hardcoded regions in `assume_role` ARNs

**Example violations:**
```hcl
provider "aws" {
  region = "us-east-1"  # ❌ Hardcoded region
}
```

**Recommended fix:**
```hcl
provider "aws" {
  region = var.aws_region  # ✅ Use variables
}
```




## aws_service_principal_hardcoded

This rule checks for hardcoded AWS service principal DNS suffixes in IAM policies and roles. Service principal DNS suffixes vary by partition (e.g., `amazonaws.com`, `amazonaws.com.cn`, `amazonaws-us-gov.com`), so hardcoding them prevents multi-partition compatibility.

**This rule is disabled by default** - enable it in your `.tflint.hcl` if you want to enforce this check.

**Example violations:**
```hcl
resource "aws_iam_role" "bad" {
  assume_role_policy = jsonencode({
    Statement = [{
      Principal = {
        Service = "s3.amazonaws.com"  # ❌ Hardcoded DNS suffix
      }
    }]
  })
}

resource "aws_iam_role" "bad_china" {
  assume_role_policy = jsonencode({
    Statement = [{
      Principal = {
        Service = "lambda.amazonaws.com.cn"  # ❌ Hardcoded China DNS suffix
      }
    }]
  })
}

resource "aws_iam_role" "bad_multiple" {
  assume_role_policy = jsonencode({
    Statement = [{
      Principal = {
        Service = [
          "lambda.amazonaws.com",      # ❌ Hardcoded
          "ec2.amazonaws.com",          # ❌ Hardcoded
          "ecs-tasks.amazonaws.com"     # ❌ Hardcoded
        ]
      }
    }]
  })
}
```

**Recommended fix:**
```hcl
data "aws_service_principal" "s3" {
  service_name = "s3"
}

resource "aws_iam_role" "good" {
  assume_role_policy = jsonencode({
    Statement = [{
      Principal = {
        Service = data.aws_service_principal.s3.name  # ✅ Dynamic
      }
    }]
  })
}

data "aws_service_principal" "lambda" {
  service_name = "lambda"
}

data "aws_service_principal" "ec2" {
  service_name = "ec2"
}

data "aws_service_principal" "ecs_tasks" {
  service_name = "ecs-tasks"
}

resource "aws_iam_role" "good_multiple" {
  assume_role_policy = jsonencode({
    Statement = [{
      Principal = {
        Service = [
          data.aws_service_principal.lambda.name,      # ✅ Dynamic
          data.aws_service_principal.ec2.name,          # ✅ Dynamic
          data.aws_service_principal.ecs_tasks.name     # ✅ Dynamic
        ]
      }
    }]
  })
}
```

## aws_service_principal_dns_suffix

This rule checks for use of `dns_suffix` in service principal construction. While using `data.aws_partition.current.dns_suffix` is better than hardcoding, the `aws_service_principal` data source is the recommended best practice as it handles all the complexity for you.

**Example violations:**
```hcl
data "aws_partition" "current" {}

resource "aws_iam_role" "bad" {
  assume_role_policy = jsonencode({
    Statement = [{
      Principal = {
        Service = "s3.${data.aws_partition.current.dns_suffix}"  # ❌ Using dns_suffix
      }
    }]
  })
}

resource "aws_iam_role" "bad_multiple" {
  assume_role_policy = jsonencode({
    Statement = [{
      Principal = {
        Service = [
          "lambda.${data.aws_partition.current.dns_suffix}",      # ❌ Using dns_suffix
          "ec2.${data.aws_partition.current.dns_suffix}",          # ❌ Using dns_suffix
          "ecs-tasks.${data.aws_partition.current.dns_suffix}"     # ❌ Using dns_suffix
        ]
      }
    }]
  })
}
```

**Recommended fix:**
```hcl
data "aws_service_principal" "s3" {
  service_name = "s3"
}

resource "aws_iam_role" "good" {
  assume_role_policy = jsonencode({
    Statement = [{
      Principal = {
        Service = data.aws_service_principal.s3.name  # ✅ Best practice
      }
    }]
  })
}

data "aws_service_principal" "lambda" {
  service_name = "lambda"
}

data "aws_service_principal" "ec2" {
  service_name = "ec2"
}

data "aws_service_principal" "ecs_tasks" {
  service_name = "ecs-tasks"
}

resource "aws_iam_role" "good_multiple" {
  assume_role_policy = jsonencode({
    Statement = [{
      Principal = {
        Service = [
          data.aws_service_principal.lambda.name,      # ✅ Best practice
          data.aws_service_principal.ec2.name,          # ✅ Best practice
          data.aws_service_principal.ecs_tasks.name     # ✅ Best practice
        ]
      }
    }]
  })
}
```
