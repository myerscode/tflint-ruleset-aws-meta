package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/myerscode/tflint-ruleset-aws-meta/rules/awsmeta"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsProviderHardcodedRegionRule checks for hardcoded AWS regions in provider configuration
type AwsProviderHardcodedRegionRule struct {
	tflint.DefaultRule
}

// NewAwsProviderHardcodedRegionRule returns a new rule
func NewAwsProviderHardcodedRegionRule() *AwsProviderHardcodedRegionRule {
	return &AwsProviderHardcodedRegionRule{}
}

// Name returns the rule name
func (r *AwsProviderHardcodedRegionRule) Name() string {
	return "aws_provider_hardcoded_region"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsProviderHardcodedRegionRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *AwsProviderHardcodedRegionRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsProviderHardcodedRegionRule) Link() string {
	return ""
}

// Check checks for hardcoded AWS regions in provider configuration
func (r *AwsProviderHardcodedRegionRule) Check(runner tflint.Runner) error {
	regionPattern := awsmeta.GetRegionPattern()
	arnRegionPattern := awsmeta.GetARNRegionPattern()

	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	providers, err := runner.GetProviderContent("aws", &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: "region"},
		},
		Blocks: []hclext.BlockSchema{
			{
				Type: "assume_role",
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "role_arn"},
					},
				},
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, provider := range providers.Blocks {
		if attr, exists := provider.Body.Attributes["region"]; exists {
			// Skip if the expression is not a literal (e.g. var.region, local.region)
			if isReference(files, attr.Expr.Range()) {
				continue
			}

			err := runner.EvaluateExpr(attr.Expr, func(region string) error {
				if regionPattern.MatchString(region) {
					return runner.EmitIssue(
						r,
						fmt.Sprintf("Hardcoded AWS region '%s' in provider configuration. Consider using variables or environment variables for better flexibility", region),
						attr.Expr.Range(),
					)
				}
				return nil
			}, nil)
			if err != nil && !strings.Contains(err.Error(), "cannot convert") {
				return err
			}
		}

		for _, assumeRoleBlock := range provider.Body.Blocks {
			if assumeRoleBlock.Type == "assume_role" {
				if attr, exists := assumeRoleBlock.Body.Attributes["role_arn"]; exists {
					if isReference(files, attr.Expr.Range()) {
						continue
					}

					err := runner.EvaluateExpr(attr.Expr, func(roleArn string) error {
						if matches := arnRegionPattern.FindStringSubmatch(roleArn); len(matches) > 1 {
							region := matches[1]
							return runner.EmitIssue(
								r,
								fmt.Sprintf("Hardcoded AWS region '%s' found in assume_role ARN. Consider using variables or data.aws_region.current.name", region),
								attr.Expr.Range(),
							)
						}
						return nil
					}, nil)
					if err != nil && !strings.Contains(err.Error(), "cannot convert") {
						return err
					}
				}
			}
		}
	}

	return nil
}

// isReference checks if the expression source text contains variable, local, or data references
func isReference(files map[string]*hcl.File, exprRange hcl.Range) bool {
	if file, ok := files[exprRange.Filename]; ok {
		src := file.Bytes
		if exprRange.Start.Byte < len(src) && exprRange.End.Byte <= len(src) {
			sourceText := string(src[exprRange.Start.Byte:exprRange.End.Byte])
			return strings.Contains(sourceText, "var.") ||
				strings.Contains(sourceText, "local.") ||
				strings.Contains(sourceText, "data.") ||
				strings.Contains(sourceText, "module.")
		}
	}
	return false
}
