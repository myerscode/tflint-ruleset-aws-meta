# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.5.1] - 2026-04-02
### :bug: Bug Fixes
- [`628bebe`](https://github.com/myerscode/tflint-ruleset-aws-meta/commit/628bebe4c9618f5bcfd9f7ecdd84ddf36813b941) - **ci**: prevent duplicate releases *(commit by [@oniice](https://github.com/oniice))*
- [`a29226a`](https://github.com/myerscode/tflint-ruleset-aws-meta/commit/a29226a9d3981a38ef0024097907a72f82556dd6) - **rules**: skip variable references in provider region check *(commit by [@oniice](https://github.com/oniice))*


## [v0.5.0] - 2026-04-02
### :sparkles: New Features
- [`85e42f3`](https://github.com/myerscode/tflint-ruleset-aws-meta/commit/85e42f3ffbcdaa5eaaaf615cc4c5759e91afc908) - **release**: inject version from git tag *(commit by [@oniice](https://github.com/oniice))*
- [`eb8430a`](https://github.com/myerscode/tflint-ruleset-aws-meta/commit/eb8430a5b391e29464f728c1defca2209021da97) - **rules**: detect hardcoded AZs and regions *(commit by [@oniice](https://github.com/oniice))*
- [`253a450`](https://github.com/myerscode/tflint-ruleset-aws-meta/commit/253a45090b83395012289b96338edf486c4d9292) - **rules**: add aws_hardcoded_ids rule *(commit by [@oniice](https://github.com/oniice))*
- [`299357a`](https://github.com/myerscode/tflint-ruleset-aws-meta/commit/299357a1417c0c946dc32e1e3e185620bde12537) - **ci**: add manual dispatch with semver validation *(commit by [@oniice](https://github.com/oniice))*

### :wrench: Chores
- [`a18624a`](https://github.com/myerscode/tflint-ruleset-aws-meta/commit/a18624aa5da20db9b64f67d6dfa1c12376053c68) - **deps**: Bump github.com/terraform-linters/tflint-plugin-sdk *(PR [#33](https://github.com/myerscode/tflint-ruleset-aws-meta/pull/33) by [@dependabot[bot]](https://github.com/apps/dependabot))*
- [`73e9e89`](https://github.com/myerscode/tflint-ruleset-aws-meta/commit/73e9e8990285eb25d45d9d5708d486d9d21cb3ac) - **deps**: Bump github.com/myerscode/aws-meta from 0.86.0 to 0.88.0 *(PR [#37](https://github.com/myerscode/tflint-ruleset-aws-meta/pull/37) by [@dependabot[bot]](https://github.com/apps/dependabot))*
- [`00a8212`](https://github.com/myerscode/tflint-ruleset-aws-meta/commit/00a821234ea84a3d63416f60b956123a0d1b5cb2) - **deps**: Bump actions/deploy-pages from 4 to 5 *(PR [#35](https://github.com/myerscode/tflint-ruleset-aws-meta/pull/35) by [@dependabot[bot]](https://github.com/apps/dependabot))*
- [`f994ce8`](https://github.com/myerscode/tflint-ruleset-aws-meta/commit/f994ce83be0a63be04d939234595dd0cb8530c42) - **deps**: Bump actions/setup-go from 6.3.0 to 6.4.0 *(PR [#36](https://github.com/myerscode/tflint-ruleset-aws-meta/pull/36) by [@dependabot[bot]](https://github.com/apps/dependabot))*
- [`5bad0c5`](https://github.com/myerscode/tflint-ruleset-aws-meta/commit/5bad0c5b5c28c1260f51e478ede5ff3bcb7761f5) - **release**: use changelog-action for releases *(commit by [@oniice](https://github.com/oniice))*
- [`f7b8c1b`](https://github.com/myerscode/tflint-ruleset-aws-meta/commit/f7b8c1b769623896430d06d84973444872f2fff4) - **ci**: use version tags and PAT for actions *(commit by [@oniice](https://github.com/oniice))*


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

[0.4.1]: https://github.com/myerscode/tflint-ruleset-aws-meta/compare/v0.4.0...v0.4.1
[0.4.0]: https://github.com/myerscode/tflint-ruleset-aws-meta/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/myerscode/tflint-ruleset-aws-meta/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/myerscode/tflint-ruleset-aws-meta/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/myerscode/tflint-ruleset-aws-meta/releases/tag/v0.1.0
[v0.5.0]: https://github.com/myerscode/tflint-ruleset-aws-meta/compare/v0.4.1...v0.5.0
[v0.5.1]: https://github.com/myerscode/tflint-ruleset-aws-meta/compare/v0.5.0...v0.5.1
