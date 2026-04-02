package rules

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/myerscode/tflint-ruleset-aws-meta/rules/awsmeta"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsHardcodedIDsRule checks for hardcoded AWS account IDs and AMI IDs
type AwsHardcodedIDsRule struct {
	tflint.DefaultRule
}

// NewAwsHardcodedIDsRule returns a new rule
func NewAwsHardcodedIDsRule() *AwsHardcodedIDsRule {
	return &AwsHardcodedIDsRule{}
}

// Name returns the rule name
func (r *AwsHardcodedIDsRule) Name() string {
	return "aws_hardcoded_ids"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsHardcodedIDsRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *AwsHardcodedIDsRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsHardcodedIDsRule) Link() string {
	return ""
}

// Check checks for hardcoded AWS account IDs and AMI IDs
func (r *AwsHardcodedIDsRule) Check(runner tflint.Runner) error {
	accountIDPattern := awsmeta.GetAccountIDPattern()
	amiIDPattern := awsmeta.GetAMIIDPattern()

	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	checked := make(map[string]bool)

	diags := runner.WalkExpressions(tflint.ExprWalkFunc(func(expr hcl.Expression) hcl.Diagnostics {
		exprKey := fmt.Sprintf("%s:%d:%d", expr.Range().Filename, expr.Range().Start.Line, expr.Range().Start.Column)
		if checked[exprKey] {
			return nil
		}
		checked[exprKey] = true

		// Pre-filter on raw source text
		exprRange := expr.Range()
		if file, ok := files[exprRange.Filename]; ok {
			src := file.Bytes
			if exprRange.Start.Byte < len(src) && exprRange.End.Byte <= len(src) {
				sourceText := string(src[exprRange.Start.Byte:exprRange.End.Byte])
				hasAccountID := accountIDPattern.MatchString(sourceText)
				hasAMI := amiIDPattern.MatchString(sourceText)
				if !hasAccountID && !hasAMI {
					return nil
				}
			}
		}

		err := runner.EvaluateExpr(expr, func(value string) error {
			// Check for hardcoded account ID
			if matches := accountIDPattern.FindStringSubmatch(value); len(matches) > 1 {
				if err := runner.EmitIssue(
					r,
					fmt.Sprintf("Hardcoded AWS account ID '%s' found. Consider using data.aws_caller_identity.current.account_id", matches[1]),
					expr.Range(),
				); err != nil {
					return err
				}
			}

			// Check for hardcoded AMI ID
			if match := amiIDPattern.FindString(value); match != "" {
				if err := runner.EmitIssue(
					r,
					fmt.Sprintf("Hardcoded AMI ID '%s' found. AMI IDs are region-specific. Consider using data.aws_ami to dynamically look up AMIs", match),
					expr.Range(),
				); err != nil {
					return err
				}
			}

			return nil
		}, nil)

		_ = err

		return nil
	}))

	if diags.HasErrors() {
		return diags
	}

	return nil
}
