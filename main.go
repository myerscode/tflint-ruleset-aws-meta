package main

import (
	"github.com/myerscode/tflint-ruleset-aws-meta/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "aws-meta",
			Version: "0.1.0",
			Rules: []tflint.Rule{
				rules.NewAwsIamRolePolicyHardcodedRegionRule(),
				rules.NewAwsIamRolePolicyHardcodedPartitionRule(),
				rules.NewAwsIamPolicyHardcodedRegionRule(),
				rules.NewAwsIamPolicyHardcodedPartitionRule(),
				rules.NewAwsProviderHardcodedRegionRule(),
				rules.NewAwsARNHardcodedRule(),
			},
		},
	})
}
