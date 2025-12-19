---
title: aws_service_principal_hardcoded
description: Detects hardcoded AWS service principal DNS suffixes.
---

This rule checks for hardcoded AWS service principal DNS suffixes in expressions and strings across Terraform files.

It matches service principal forms like `service.amazonaws.com`, `service.amazonaws.com.cn`, and `service.amazonaws-us-gov.com` and emits an issue suggesting use of data sources (e.g., `data.aws_service_principal.<name>.name`) for multi-partition compatibility.

Example: a string like `s3.amazonaws.com` embedded in a policy or principal could be flagged as hardcoded.
