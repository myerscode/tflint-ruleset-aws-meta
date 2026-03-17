# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.4.1] - 2026-03-17

### Changed

- Enable `aws_service_principal_hardcoded` rule by default

## [0.4.0] - 2026-03-16

### Added

- Dynamic DNS suffix pattern for service principal rules
- Regex caching with `sync.Once` for all pattern functions
- Source-text pre-filter for `aws_meta_hardcoded` and `aws_service_principal_dns_suffix` rules
- Documentation site published to GitHub Pages

### Changed

- Bumped `aws-meta` from 0.70.0 to 0.86.0
- Bumped `tflint-plugin-sdk` to 0.23.1
- Bumped GitHub Actions dependencies (checkout, setup-go, setup-node, goreleaser-action, attest-build-provenance, upload-pages-artifact)

### Fixed

- Made documentation links relative to rules directory

### Removed

- MkDocs configuration (replaced by Astro/Starlight docs site)

## [0.3.0] - 2025-11-26

### Added

- `aws_service_principal_hardcoded` rule for detecting hardcoded DNS suffixes in service principals
- `aws_service_principal_dns_suffix` rule for detecting dns_suffix interpolation in service principals

### Changed

- Bumped `aws-meta` from 0.69.0 to 0.70.0
- Bumped GitHub Actions dependencies (checkout, setup-go)

## [0.2.0] - 2025-11-21

### Changed

- Renamed `aws_hardcoded_region` rule to `aws_provider_hardcoded_region`
- Made `aws_provider_hardcoded_region` rule optional (disabled by default)

## [0.1.0] - 2025-11-20

### Added

- Initial release
- `aws_meta_hardcoded` rule for detecting hardcoded regions and partitions in ARN values across all AWS resource types
- `aws_iam_policy_hardcoded_region` rule
- `aws_iam_policy_hardcoded_partition` rule
- `aws_iam_role_policy_hardcoded_region` rule
- `aws_iam_role_policy_hardcoded_partition` rule
- `aws_provider_hardcoded_region` rule
- Passing and failing example configurations
- Dynamic pattern generation using `aws-meta` package

[Unreleased]: https://github.com/myerscode/tflint-ruleset-aws-meta/compare/v0.4.0...HEAD
[0.4.0]: https://github.com/myerscode/tflint-ruleset-aws-meta/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/myerscode/tflint-ruleset-aws-meta/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/myerscode/tflint-ruleset-aws-meta/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/myerscode/tflint-ruleset-aws-meta/releases/tag/v0.1.0
