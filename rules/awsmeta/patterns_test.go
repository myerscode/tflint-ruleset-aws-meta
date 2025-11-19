package awsmeta

import (
	"testing"
)

func TestGetRegionPattern(t *testing.T) {
	pattern := GetRegionPattern()

	// Test some known regions
	testCases := []struct {
		region   string
		expected bool
	}{
		{"us-east-1", true},
		{"us-west-2", true},
		{"eu-west-1", true},
		{"ap-southeast-1", true},
		{"invalid-region", false},
		{"", false},
	}

	for _, tc := range testCases {
		result := pattern.MatchString(tc.region)
		if result != tc.expected {
			t.Errorf("Region %s: expected %v, got %v", tc.region, tc.expected, result)
		}
	}
}

func TestGetAvailabilityZonePattern(t *testing.T) {
	pattern := GetAvailabilityZonePattern()

	// Test some known AZs
	testCases := []struct {
		az       string
		expected bool
	}{
		{"us-east-1a", true},
		{"us-west-2b", true},
		{"eu-west-1c", true},
		{"us-east-1", false}, // Not an AZ, just a region
		{"invalid-az", false},
	}

	for _, tc := range testCases {
		result := pattern.MatchString(tc.az)
		if result != tc.expected {
			t.Errorf("AZ %s: expected %v, got %v", tc.az, tc.expected, result)
		}
	}
}

func TestGetARNRegionPattern(t *testing.T) {
	pattern := GetARNRegionPattern()

	// Test ARNs with regions
	testCases := []struct {
		arn      string
		expected bool
	}{
		{"arn:aws:s3:us-east-1:123456789012:bucket/my-bucket", true},
		{"arn:aws:ec2:eu-west-1:123456789012:instance/i-1234567890abcdef0", true},
		{"arn:aws:iam::123456789012:role/my-role", false}, // No region in IAM ARN
		{"not-an-arn", false},
	}

	for _, tc := range testCases {
		result := pattern.MatchString(tc.arn)
		if result != tc.expected {
			t.Errorf("ARN %s: expected %v, got %v", tc.arn, tc.expected, result)
		}
	}
}

func TestGetPartitionPattern(t *testing.T) {
	pattern := GetPartitionPattern()

	// Test ARNs with different partitions
	testCases := []struct {
		arn      string
		expected bool
	}{
		{"arn:aws:s3:us-east-1:123456789012:bucket/my-bucket", true},
		{"arn:aws-cn:s3:cn-north-1:123456789012:bucket/my-bucket", true},
		{"arn:aws-us-gov:s3:us-gov-west-1:123456789012:bucket/my-bucket", true},
		{"not-an-arn", false},
	}

	for _, tc := range testCases {
		result := pattern.MatchString(tc.arn)
		if result != tc.expected {
			t.Errorf("ARN %s: expected %v, got %v", tc.arn, tc.expected, result)
		}
	}
}

func TestGetRegionInStringPattern(t *testing.T) {
	pattern := GetRegionInStringPattern()

	// Test finding regions within strings
	testCases := []struct {
		text     string
		expected bool
	}{
		{"The bucket is in us-east-1", true},
		{"Deploying to eu-west-1 region", true},
		{"No region here", false},
		{"", false},
	}

	for _, tc := range testCases {
		result := pattern.MatchString(tc.text)
		if result != tc.expected {
			t.Errorf("Text %q: expected %v, got %v", tc.text, tc.expected, result)
		}
	}
}

