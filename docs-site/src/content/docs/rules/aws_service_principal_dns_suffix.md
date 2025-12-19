---
title: aws_service_principal_dns_suffix
description: Detects use of `dns_suffix` interpolation in service principals.
---

This rule checks for use of `dns_suffix` in service principals (e.g., `service.${var.dns_suffix}`) and suggests using `data.aws_service_principal.<name>.name` instead for better maintainability.

It detects both evaluated strings containing `dns_suffix` and raw interpolated expressions that reference `.dns_suffix`.
