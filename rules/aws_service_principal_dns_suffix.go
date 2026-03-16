package rules

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsServicePrincipalDNSSuffixRule checks for use of dns_suffix in service principals
type AwsServicePrincipalDNSSuffixRule struct {
	tflint.DefaultRule
}

// NewAwsServicePrincipalDNSSuffixRule returns a new rule
func NewAwsServicePrincipalDNSSuffixRule() *AwsServicePrincipalDNSSuffixRule {
	return &AwsServicePrincipalDNSSuffixRule{}
}

// Name returns the rule name
func (r *AwsServicePrincipalDNSSuffixRule) Name() string {
	return "aws_service_principal_dns_suffix"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsServicePrincipalDNSSuffixRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsServicePrincipalDNSSuffixRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsServicePrincipalDNSSuffixRule) Link() string {
	return ""
}

// Pattern to match service names with dns_suffix interpolation
var dnsSuffixPattern = regexp.MustCompile(`([a-z0-9\-]+)\.\$\{[^}]*\.dns_suffix\}`)

// Check checks for use of dns_suffix in service principals
func (r *AwsServicePrincipalDNSSuffixRule) Check(runner tflint.Runner) error {
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

		// Pre-filter: check raw source for "dns_suffix" before making gRPC call
		exprRange := expr.Range()
		var sourceText string
		if file, ok := files[exprRange.Filename]; ok {
			src := file.Bytes
			if exprRange.Start.Byte < len(src) && exprRange.End.Byte <= len(src) {
				sourceText = string(src[exprRange.Start.Byte:exprRange.End.Byte])
				if !strings.Contains(sourceText, "dns_suffix") {
					return nil
				}
			}
		}

		// Skip pure variable/attribute references (e.g. data.aws_partition.current.dns_suffix)
		// These contain "dns_suffix" in their source but aren't hardcoded strings
		if _, diags := hcl.AbsTraversalForExpr(expr); !diags.HasErrors() {
			return nil
		}

		// Try to evaluate the expression as a string
		err := runner.EvaluateExpr(expr, func(value string) error {
			if strings.Contains(value, "dns_suffix") {
				if matches := dnsSuffixPattern.FindStringSubmatch(value); len(matches) > 1 {
					serviceName := matches[1]
					if err := runner.EmitIssue(
						r,
						fmt.Sprintf("Service principal uses dns_suffix. Consider using data.aws_service_principal.%s.name instead for better maintainability", strings.ReplaceAll(serviceName, "-", "_")),
						expr.Range(),
					); err != nil {
						return err
					}
				}
			}
			return nil
		}, nil)

		// If evaluation failed, check the raw source text directly
		if err != nil && strings.Contains(sourceText, "dns_suffix") {
			if matches := regexp.MustCompile(`"([a-z0-9\-]+)\.\$\{[^}]*\.dns_suffix\}"`).FindStringSubmatch(sourceText); len(matches) > 1 {
				serviceName := matches[1]
				_ = runner.EmitIssue(
					r,
					fmt.Sprintf("Service principal uses dns_suffix. Consider using data.aws_service_principal.%s.name instead for better maintainability", strings.ReplaceAll(serviceName, "-", "_")),
					expr.Range(),
				)
			} else {
				_ = runner.EmitIssue(
					r,
					"Service principal uses dns_suffix. Consider using data.aws_service_principal data source instead for better maintainability",
					expr.Range(),
				)
			}
		}

		return nil
	}))

	if diags.HasErrors() {
		return diags
	}

	return nil
}
