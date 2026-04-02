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
// Check checks for hardcoded regions and partitions in ARN-like string values
func (r *AwsMetaHardcodedRule) Check(runner tflint.Runner) error {
	arnRegionPattern := awsmeta.GetARNRegionPattern()
	arnPartitionPattern := awsmeta.GetPartitionPattern()
	azPattern := awsmeta.GetAvailabilityZonePattern()
	regionPattern := awsmeta.GetRegionPattern()

	// Get all source files upfront so we can inspect raw expression text
	// before making expensive gRPC EvaluateExpr calls
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	// Track which expressions we've already checked to avoid duplicates
	checked := make(map[string]bool)

	// Walk all expressions in the Terraform files
	diags := runner.WalkExpressions(tflint.ExprWalkFunc(func(expr hcl.Expression) hcl.Diagnostics {
		exprKey := fmt.Sprintf("%s:%d:%d", expr.Range().Filename, expr.Range().Start.Line, expr.Range().Start.Column)
		if checked[exprKey] {
			return nil
		}
		checked[exprKey] = true

		// Pre-filter: check the raw source text for potential matches before
		// making the expensive gRPC EvaluateExpr call.
		exprRange := expr.Range()
		if file, ok := files[exprRange.Filename]; ok {
			src := file.Bytes
			if exprRange.Start.Byte < len(src) && exprRange.End.Byte <= len(src) {
				sourceText := strings.ToLower(string(src[exprRange.Start.Byte:exprRange.End.Byte]))
				hasARN := strings.Contains(sourceText, "arn:")
				hasRegionLike := awsmeta.GetRegionInStringPattern().MatchString(sourceText)
				if !hasARN && !hasRegionLike {
					return nil
				}
			}
		}

		// Try to evaluate the expression as a string
		err := runner.EvaluateExpr(expr, func(value string) error {
			// Check for hardcoded ARN values
			if strings.HasPrefix(value, "arn:") {
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
			}

			// Check for hardcoded availability zone (e.g. "eu-west-2a")
			if azPattern.MatchString(value) {
				// Extract the region part (everything except the last character)
				region := value[:len(value)-1]
				if err := runner.EmitIssue(
					r,
					fmt.Sprintf("Hardcoded AWS availability zone '%s' found. Consider using data.aws_availability_zones to dynamically fetch AZs for the current region", value),
					expr.Range(),
				); err != nil {
					return err
				}
				_ = region
				return nil
			}

			// Check for hardcoded region as a standalone value (e.g. "eu-west-2")
			if regionPattern.MatchString(value) {
				if err := runner.EmitIssue(
					r,
					fmt.Sprintf("Hardcoded AWS region '%s' found. Consider using data.aws_region.current.name", value),
					expr.Range(),
				); err != nil {
					return err
				}
			}

			return nil
		}, nil)

		// Silently ignore evaluation errors
		_ = err

		return nil
	}))

	if diags.HasErrors() {
		return diags
	}

	return nil
}
