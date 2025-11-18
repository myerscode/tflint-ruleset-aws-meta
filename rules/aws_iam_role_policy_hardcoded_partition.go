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

// AwsIamRolePolicyHardcodedPartitionRule checks for hardcoded AWS partitions in IAM role policies
type AwsIamRolePolicyHardcodedPartitionRule struct {
	tflint.DefaultRule
}

// NewAwsIamRolePolicyHardcodedPartitionRule returns a new rule
func NewAwsIamRolePolicyHardcodedPartitionRule() *AwsIamRolePolicyHardcodedPartitionRule {
	return &AwsIamRolePolicyHardcodedPartitionRule{}
}

// Name returns the rule name
func (r *AwsIamRolePolicyHardcodedPartitionRule) Name() string {
	return "aws_iam_role_policy_hardcoded_partition"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsIamRolePolicyHardcodedPartitionRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsIamRolePolicyHardcodedPartitionRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsIamRolePolicyHardcodedPartitionRule) Link() string {
	return ""
}

// Check checks for hardcoded AWS partitions in IAM role policies
func (r *AwsIamRolePolicyHardcodedPartitionRule) Check(runner tflint.Runner) error {
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
				return r.checkPolicyForHardcodedPartitions(runner, policy, attr.Expr.Range())
			}, nil)
			if err != nil && !strings.Contains(err.Error(), "cannot convert") {
				return err
			}
		}
	}

	return nil
}

func (r *AwsIamRolePolicyHardcodedPartitionRule) checkPolicyForHardcodedPartitions(runner tflint.Runner, policy string, rng hcl.Range) error {
	// Get dynamic pattern from aws-meta package
	arnPartitionPattern := awsmeta.GetPartitionPattern()

	// Try to parse as JSON to check structured policy
	var policyDoc map[string]interface{}
	if err := json.Unmarshal([]byte(policy), &policyDoc); err == nil {
		// Check structured policy document
		if err := r.checkPolicyDocument(runner, policyDoc, rng, arnPartitionPattern); err != nil {
			return err
		}
	} else {
		// Check raw string for ARN patterns
		if matches := arnPartitionPattern.FindAllStringSubmatch(policy, -1); len(matches) > 0 {
			for _, match := range matches {
				if len(match) > 1 {
					partition := match[1]
					if err := runner.EmitIssue(
						r,
						fmt.Sprintf("Hardcoded AWS partition '%s' found in ARN within IAM role policy. Consider using data.aws_partition.current.partition", partition),
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

func (r *AwsIamRolePolicyHardcodedPartitionRule) checkPolicyDocument(runner tflint.Runner, doc map[string]interface{}, rng hcl.Range, arnPartitionPattern *regexp.Regexp) error {
	// Convert back to string to search for patterns
	docBytes, err := json.Marshal(doc)
	if err != nil {
		return nil // Skip if we can't marshal back
	}

	docString := string(docBytes)

	// Check for hardcoded partitions in ARNs within the policy document
	if matches := arnPartitionPattern.FindAllStringSubmatch(docString, -1); len(matches) > 0 {
		for _, match := range matches {
			if len(match) > 1 {
				partition := match[1]
				if err := runner.EmitIssue(
					r,
					fmt.Sprintf("Hardcoded AWS partition '%s' found in ARN within IAM role policy document. Consider using data.aws_partition.current.partition", partition),
					rng,
				); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
