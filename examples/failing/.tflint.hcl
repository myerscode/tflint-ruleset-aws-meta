plugin "aws-meta" {
  enabled = true
}

rule "aws_iam_role_policy_hardcoded_region" {
  enabled = true
}

rule "aws_iam_role_policy_hardcoded_partition" {
  enabled = true
}

rule "aws_iam_policy_hardcoded_region" {
  enabled = true
}

rule "aws_iam_policy_hardcoded_partition" {
  enabled = true
}

rule "aws_provider_hardcoded_region" {
  enabled = true
}

rule "aws_arn_hardcoded" {
  enabled = true
}