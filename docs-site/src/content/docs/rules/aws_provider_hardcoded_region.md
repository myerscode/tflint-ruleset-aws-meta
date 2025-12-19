---
title: aws_provider_hardcoded_region
description: Checks AWS provider configurations for hardcoded regions.
---

This rule checks AWS provider configurations for hardcoded regions. It detects:

- Hardcoded regions in provider `region` attribute
- Hardcoded regions in `assume_role` ARNs

## Example violations

```hcl
provider "aws" {
  region = "us-east-1"  # ❌ Hardcoded region
}
```

## Recommended fix

```hcl
provider "aws" {
  region = var.aws_region  # ✅ Use variables
}
```
