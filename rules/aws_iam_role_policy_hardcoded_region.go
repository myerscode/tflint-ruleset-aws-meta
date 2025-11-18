package rules

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/myerscode/tflint-ruleset-aws-meta/rules/awsmeta"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsIamRolePolicyHardcodedRegionRule checks for hardcoded AWS regions in IAM role policies
type AwsIamRolePolicyHardcodedRegionRule struct {
	tflint.DefaultRule
}

// NewAwsIamRolePolicyHardcodedRegionRule returns a new rule
func NewAwsIamRolePolicyHardcodedRegionRule() *AwsIamRolePolicyHardcodedRegionRule {
	return &AwsIamRolePolicyHardcodedRegionRule{}
}

// Name returns the rule name
func (r *AwsIamRolePolicyHardcodedRegionRule) Name() string {
	return "aws_iam_role_policy_hardcoded_region"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsIamRolePolicyHardcodedRegionRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsIamRolePolicyHardcodedRegionRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsIamRolePolicyHardcodedRegionRule) Link() string {
	return ""
}

// Check checks for hardcoded AWS regions in IAM role policies
func (r *AwsIamRolePolicyHardcodedRegionRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent("aws_iam_role_policy", &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: "policy"},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		if attr, exists := resource.Body.Attributes["policy"]; exists {
			err := runner.EvaluateExpr(attr.Expr, func(policy string) error {
				return r.checkPolicyForHardcodedRegions(runner, policy, attr.Expr.Range())
			}, nil)
			if err != nil && !strings.Contains(err.Error(), "cannot convert") {
				return err
			}
		}
	}

	return nil
}

func (r *AwsIamRolePolicyHardcodedRegionRule) checkPolicyForHardcodedRegions(runner tflint.Runner, policy string, rng hcl.Range) error {
	// Get dynamic patterns from aws-meta package
	regionInStringPattern := awsmeta.GetRegionInStringPattern()
	arnRegionPattern := awsmeta.GetARNRegionPattern()

	// Try to parse as JSON to check structured policy
	var policyDoc map[string]interface{}
	if err := json.Unmarshal([]byte(policy), &policyDoc); err == nil {
		// Check structured policy document
		if err := r.checkPolicyDocument(runner, policyDoc, rng, regionInStringPattern, arnRegionPattern); err != nil {
			return err
		}
	} else {
		// Check raw string for patterns
		if matches := regionInStringPattern.FindAllString(policy, -1); len(matches) > 0 {
			for _, match := range matches {
				if err := runner.EmitIssue(
					r,
					fmt.Sprintf("Hardcoded AWS region '%s' found in IAM role policy. Consider using variables or data.aws_region.current.name", match),
					rng,
				); err != nil {
					return err
				}
			}
		}

		if matches := arnRegionPattern.FindAllStringSubmatch(policy, -1); len(matches) > 0 {
			for _, match := range matches {
				if len(match) > 1 {
					region := match[1]
					if err := runner.EmitIssue(
						r,
						fmt.Sprintf("Hardcoded AWS region '%s' found in ARN within IAM role policy. Consider using variables or data.aws_region.current.name", region),
						rng,
					); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func (r *AwsIamRolePolicyHardcodedRegionRule) checkPolicyDocument(runner tflint.Runner, doc map[string]interface{}, rng hcl.Range, regionInStringPattern, arnRegionPattern *regexp.Regexp) error {
	// Convert back to string to search for patterns
	docBytes, err := json.Marshal(doc)
	if err != nil {
		return nil // Skip if we can't marshal back
	}

	docString := string(docBytes)

	// Check for hardcoded regions in the policy document
	if matches := regionInStringPattern.FindAllString(docString, -1); len(matches) > 0 {
		for _, match := range matches {
			if err := runner.EmitIssue(
				r,
				fmt.Sprintf("Hardcoded AWS region '%s' found in IAM role policy document. Consider using variables or data.aws_region.current.name", match),
				rng,
			); err != nil {
				return err
			}
		}
	}

	if matches := arnRegionPattern.FindAllStringSubmatch(docString, -1); len(matches) > 0 {
		for _, match := range matches {
			if len(match) > 1 {
				region := match[1]
				if err := runner.EmitIssue(
					r,
					fmt.Sprintf("Hardcoded AWS region '%s' found in ARN within IAM role policy document. Consider using variables or data.aws_region.current.name", region),
					rng,
				); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
