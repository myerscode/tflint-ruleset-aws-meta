package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/myerscode/tflint-ruleset-aws-meta/rules/awsmeta"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsMetaHardcodedRule checks for hardcoded regions and partitions in ARN values
// across all AWS resources by walking all expressions
type AwsMetaHardcodedRule struct {
	tflint.DefaultRule
}

// NewAwsMetaHardcodedRule returns a new rule
func NewAwsMetaHardcodedRule() *AwsMetaHardcodedRule {
	return &AwsMetaHardcodedRule{}
}

// Name returns the rule name
func (r *AwsMetaHardcodedRule) Name() string {
	return "aws_meta_hardcoded"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsMetaHardcodedRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsMetaHardcodedRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsMetaHardcodedRule) Link() string {
	return ""
}

// Check checks for hardcoded regions and partitions in ARN-like string values
func (r *AwsMetaHardcodedRule) Check(runner tflint.Runner) error {
	arnRegionPattern := awsmeta.GetARNRegionPattern()
	arnPartitionPattern := awsmeta.GetPartitionPattern()

	// Track which expressions we've already checked to avoid duplicates
	checked := make(map[string]bool)

	// Walk all expressions in the Terraform files
	diags := runner.WalkExpressions(tflint.ExprWalkFunc(func(expr hcl.Expression) hcl.Diagnostics {
		// Skip if we've already checked this expression
		// Note: ExprWalkFunc is called for both Enter and Exit, so we deduplicate
		exprKey := fmt.Sprintf("%s:%d:%d", expr.Range().Filename, expr.Range().Start.Line, expr.Range().Start.Column)
		if checked[exprKey] {
			return nil
		}
		checked[exprKey] = true

		// Try to evaluate the expression as a string
		err := runner.EvaluateExpr(expr, func(value string) error {
			// Only check if it looks like an ARN
			if !strings.HasPrefix(value, "arn:") {
				return nil
			}

			// Check for hardcoded region in ARN
			if matches := arnRegionPattern.FindStringSubmatch(value); len(matches) > 1 {
				region := matches[1]
				if err := runner.EmitIssue(
					r,
					fmt.Sprintf("Hardcoded AWS region '%s' found in ARN. Consider using data.aws_region.current.name", region),
					expr.Range(),
				); err != nil {
					return err
				}
			}

			// Check for hardcoded partition in ARN
			if matches := arnPartitionPattern.FindStringSubmatch(value); len(matches) > 1 {
				partition := matches[1]
				if err := runner.EmitIssue(
					r,
					fmt.Sprintf("Hardcoded AWS partition '%s' found in ARN. Consider using data.aws_partition.current.partition", partition),
					expr.Range(),
				); err != nil {
					return err
				}
			}

			return nil
		}, nil)

		// Silently ignore evaluation errors (variables, data sources, functions, etc.)
		_ = err

		return nil
	}))

	if diags.HasErrors() {
		return diags
	}

	return nil
}
