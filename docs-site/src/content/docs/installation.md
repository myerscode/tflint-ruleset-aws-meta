---
title: Installation
description: How to install the TFLint AWS Meta Ruleset
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

## Next Steps

After installation, see the [Configuration](/configuration) guide for:
- Configuring and enabling/disabling rules
- Running TFLint with the plugin
- Configuration examples for different scenarios