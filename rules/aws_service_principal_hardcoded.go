package rules

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsServicePrincipalHardcodedRule checks for hardcoded AWS service principal DNS suffixes
type AwsServicePrincipalHardcodedRule struct {
	tflint.DefaultRule
}

// NewAwsServicePrincipalHardcodedRule returns a new rule
func NewAwsServicePrincipalHardcodedRule() *AwsServicePrincipalHardcodedRule {
	return &AwsServicePrincipalHardcodedRule{}
}

// Name returns the rule name
func (r *AwsServicePrincipalHardcodedRule) Name() string {
	return "aws_service_principal_hardcoded"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsServicePrincipalHardcodedRule) Enabled() bool {
	return false // Optional rule
}

// Severity returns the rule severity
func (r *AwsServicePrincipalHardcodedRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsServicePrincipalHardcodedRule) Link() string {
	return ""
}

// Pattern to match hardcoded service principal DNS suffixes
var servicePrincipalPattern = regexp.MustCompile(`([a-z0-9\-]+)\.(amazonaws\.com(?:\.cn)?|amazonaws-us-gov\.com)`)

// Check checks for hardcoded service principal DNS suffixes
func (r *AwsServicePrincipalHardcodedRule) Check(runner tflint.Runner) error {
	// Track which expressions we've already checked to avoid duplicates
	checked := make(map[string]bool)

	// Walk all expressions in the Terraform files
	diags := runner.WalkExpressions(tflint.ExprWalkFunc(func(expr hcl.Expression) hcl.Diagnostics {
		// Skip if we've already checked this expression
		exprKey := fmt.Sprintf("%s:%d:%d", expr.Range().Filename, expr.Range().Start.Line, expr.Range().Start.Column)
		if checked[exprKey] {
			return nil
		}
		checked[exprKey] = true

		// Try to evaluate the expression as a string
		err := runner.EvaluateExpr(expr, func(value string) error {
			// Check if it matches a hardcoded service principal pattern
			if matches := servicePrincipalPattern.FindStringSubmatch(value); len(matches) > 0 {
				serviceName := matches[1]
				dnsSuffix := matches[2]

				// Emit issue for hardcoded service principal
				if err := runner.EmitIssue(
					r,
					fmt.Sprintf("Hardcoded service principal '%s' found. Consider using data.aws_service_principal.%s.name for multi-partition compatibility", value, strings.ReplaceAll(serviceName, "-", "_")),
					expr.Range(),
				); err != nil {
					return err
				}

				_ = dnsSuffix // Keep for potential future use
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
