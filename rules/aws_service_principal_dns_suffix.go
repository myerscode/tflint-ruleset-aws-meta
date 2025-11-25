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
		// For interpolated strings, we need to check the raw expression
		err := runner.EvaluateExpr(expr, func(value string) error {
			// Check if the evaluated value contains dns_suffix pattern
			// This won't catch interpolated values, so we also check below
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

		// If evaluation failed, check the raw expression text for dns_suffix
		if err != nil {
			// Get the source text to check for dns_suffix pattern
			files, filesErr := runner.GetFiles()
			if filesErr == nil {
				if file, ok := files[expr.Range().Filename]; ok {
					sourceBytes := file.Bytes
					if expr.Range().Start.Byte < len(sourceBytes) && expr.Range().End.Byte <= len(sourceBytes) {
						sourceText := string(sourceBytes[expr.Range().Start.Byte:expr.Range().End.Byte])

						// Check if this expression contains dns_suffix
						if strings.Contains(sourceText, "dns_suffix") {
							// Try to extract service name
							if matches := regexp.MustCompile(`"([a-z0-9\-]+)\.\$\{[^}]*\.dns_suffix\}"`).FindStringSubmatch(sourceText); len(matches) > 1 {
								serviceName := matches[1]
								_ = runner.EmitIssue(
									r,
									fmt.Sprintf("Service principal uses dns_suffix. Consider using data.aws_service_principal.%s.name instead for better maintainability", strings.ReplaceAll(serviceName, "-", "_")),
									expr.Range(),
								)
							} else if strings.Contains(sourceText, "dns_suffix") {
								// Generic message if we can't extract service name
								_ = runner.EmitIssue(
									r,
									"Service principal uses dns_suffix. Consider using data.aws_service_principal data source instead for better maintainability",
									expr.Range(),
								)
							}
						}
					}
				}
			}
		}

		return nil
	}))

	if diags.HasErrors() {
		return diags
	}

	return nil
}
