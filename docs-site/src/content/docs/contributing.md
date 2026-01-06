---
title: Contributing
description: How to build, test, and contribute to the TFLint AWS Meta Ruleset
---

## Development Setup

### Prerequisites

- Go v1.25+
- Make
- Git

### Building the Plugin

Clone the repository locally and build:

```bash
git clone https://github.com/myerscode/tflint-ruleset-aws-meta.git
cd tflint-ruleset-aws-meta
make
```

### Installing for Development

You can install the built plugin locally:

```bash
make install
```

### Testing the Plugin

Run the plugin with a basic configuration:

```bash
cat << EOS > .tflint.hcl
plugin "aws-meta" {
  enabled = true
}
EOS
tflint
```

## Running Tests

### Unit Tests

Run the Go unit tests:

```bash
go test ./...
```

### Integration Tests

Test with the provided examples:

```bash
# Test passing examples (should produce no warnings)
cd examples/passing && tflint

# Test failing examples (should produce warnings)
cd examples/failing && tflint
```

## Project Structure

```
.
├── main.go                 # Plugin entry point
├── rules/                  # Rule implementations
│   ├── aws_meta_hardcoded.go
│   ├── aws_iam_*.go
│   ├── aws_provider_*.go
│   ├── aws_service_principal_*.go
│   └── awsmeta/           # Shared utilities
│       └── patterns.go    # AWS region/partition patterns
├── examples/              # Test configurations
│   ├── passing/          # Valid configurations
│   └── failing/          # Configurations with violations
├── docs-site/            # Documentation website
└── .github/workflows/    # CI/CD pipelines
```

## Adding New Rules

### 1. Create the Rule File

Create a new file in the `rules/` directory following the naming convention:

```go
// rules/aws_new_rule.go
package rules

import (
    "github.com/terraform-linters/tflint-plugin-sdk/hclext"
    "github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsNewRule checks for...
type AwsNewRule struct {
    tflint.DefaultRule
}

// NewAwsNewRule returns a new rule instance
func NewAwsNewRule() *AwsNewRule {
    return &AwsNewRule{}
}

// Name returns the rule name
func (r *AwsNewRule) Name() string {
    return "aws_new_rule"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsNewRule) Enabled() bool {
    return false // Most rules should be disabled by default
}

// Severity returns the rule severity
func (r *AwsNewRule) Severity() tflint.Severity {
    return tflint.WARNING
}

// Link returns the rule documentation URL
func (r *AwsNewRule) Link() string {
    return "https://myerscode.github.io/tflint-ruleset-aws-meta/rules/aws_new_rule"
}

// Check runs the rule logic
func (r *AwsNewRule) Check(runner tflint.Runner) error {
    // Implementation here
    return nil
}
```

### 2. Register the Rule

Add the rule to `main.go`:

```go
func main() {
    plugin.Serve(&plugin.ServeOpts{
        RuleSet: &tflint.BuiltinRuleSet{
            Name:    "aws-meta",
            Version: version,
            Rules: []tflint.Rule{
                // ... existing rules
                rules.NewAwsNewRule(),
            },
        },
    })
}
```

### 3. Add Tests

Create a test file:

```go
// rules/aws_new_rule_test.go
package rules

import (
    "testing"
    "github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsNewRule(t *testing.T) {
    cases := []struct {
        Name     string
        Content  string
        Expected helper.Issues
    }{
        {
            Name: "violation case",
            Content: `
resource "aws_example" "test" {
  // test case
}`,
            Expected: helper.Issues{
                {
                    Rule:    NewAwsNewRule(),
                    Message: "expected message",
                },
            },
        },
    }

    rule := NewAwsNewRule()
    for _, tc := range cases {
        runner := helper.TestRunner(t, map[string]string{"main.tf": tc.Content})
        if err := rule.Check(runner); err != nil {
            t.Fatalf("Unexpected error occurred: %s", err)
        }
        helper.AssertIssues(t, tc.Expected, runner.Issues)
    }
}
```

### 4. Add Documentation

Create documentation in `docs-site/src/content/docs/rules/aws_new_rule.md`:

```markdown
---
title: New Rule Description
description: Brief description of what the rule does
ruleName: aws_new_rule
---

**Rule:** `aws_new_rule`

Description of the rule...

## Example violations

```hcl
// Bad example
```

## Recommended fixes

```hcl
// Good example
```

## Enabling this rule

This rule is **disabled by default**. To enable it, add it to your `.tflint.hcl`:

```hcl
rule "aws_new_rule" {
  enabled = true
}
```
```

### 5. Update Documentation

- Add the rule to `docs-site/src/content/docs/rules/index.md`
- Add the rule to the rules table in `docs-site/src/content/docs/index.md`
- Update the main `README.md`

## Code Style

- Follow Go conventions and formatting (`gofmt`)
- Use descriptive variable and function names
- Add comments for exported functions and types
- Keep functions focused and single-purpose

## Testing Guidelines

- Write comprehensive unit tests for all rules
- Test both positive and negative cases
- Use the provided helper functions from the TFLint SDK
- Test with realistic Terraform configurations

## Submitting Changes

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests and documentation
5. Ensure all tests pass
6. Submit a pull request

## Release Process

Releases are automated through GitHub Actions when tags are pushed. The process:

1. Update version in relevant files
2. Create and push a git tag
3. GitHub Actions builds and publishes the release
4. Documentation is automatically deployed

## Getting Help

- Check existing issues and discussions
- Review the TFLint plugin SDK documentation
- Look at similar rules for implementation patterns
- Ask questions in GitHub issues