package awsmeta

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/myerscode/aws-meta/pkg/partitions"
	"github.com/myerscode/aws-meta/pkg/regions"
)

// GetRegionPattern returns a compiled regex pattern matching all AWS regions
func GetRegionPattern() *regexp.Regexp {
	regionList, err := regions.ListAllRegions()
	if err != nil {
		// Fallback to empty pattern if we can't load regions
		return regexp.MustCompile(`^$`)
	}

	regionNames := make([]string, 0, len(regionList))
	for _, region := range regionList {
		regionNames = append(regionNames, regexp.QuoteMeta(region.RegionId))
	}

	pattern := fmt.Sprintf("^(%s)$", strings.Join(regionNames, "|"))
	return regexp.MustCompile(pattern)
}

// GetRegionInStringPattern returns a compiled regex pattern for finding regions within strings
func GetRegionInStringPattern() *regexp.Regexp {
	regionList, err := regions.ListAllRegions()
	if err != nil {
		// Fallback to empty pattern if we can't load regions
		return regexp.MustCompile(`^$`)
	}

	regionNames := make([]string, 0, len(regionList))
	for _, region := range regionList {
		regionNames = append(regionNames, regexp.QuoteMeta(region.RegionId))
	}

	pattern := fmt.Sprintf("(%s)", strings.Join(regionNames, "|"))
	return regexp.MustCompile(pattern)
}

// GetAvailabilityZonePattern returns a compiled regex pattern matching availability zones
func GetAvailabilityZonePattern() *regexp.Regexp {
	regionList, err := regions.ListAllRegions()
	if err != nil {
		// Fallback to empty pattern if we can't load regions
		return regexp.MustCompile(`^$`)
	}

	regionNames := make([]string, 0, len(regionList))
	for _, region := range regionList {
		regionNames = append(regionNames, regexp.QuoteMeta(region.RegionId))
	}

	// AZ pattern: region code + letter (a-z)
	pattern := fmt.Sprintf("^(%s)[a-z]$", strings.Join(regionNames, "|"))
	return regexp.MustCompile(pattern)
}

// GetARNRegionPattern returns a compiled regex pattern for finding regions in ARNs
func GetARNRegionPattern() *regexp.Regexp {
	regionList, err := regions.ListAllRegions()
	if err != nil {
		// Fallback to empty pattern if we can't load regions
		return regexp.MustCompile(`^$`)
	}

	regionNames := make([]string, 0, len(regionList))
	for _, region := range regionList {
		regionNames = append(regionNames, regexp.QuoteMeta(region.RegionId))
	}

	pattern := fmt.Sprintf(`arn:aws[^:]*:[^:]+:(%s):`, strings.Join(regionNames, "|"))
	return regexp.MustCompile(pattern)
}

// GetPartitionPattern returns a compiled regex pattern matching all AWS partitions
func GetPartitionPattern() *regexp.Regexp {
	partitionList, err := partitions.List()
	if err != nil {
		// Fallback to empty pattern if we can't load partitions
		return regexp.MustCompile(`^$`)
	}

	partitionNames := make([]string, 0, len(partitionList))
	for _, partition := range partitionList {
		partitionNames = append(partitionNames, regexp.QuoteMeta(partition.ID))
	}

	pattern := fmt.Sprintf(`arn:(%s):`, strings.Join(partitionNames, "|"))
	return regexp.MustCompile(pattern)
}
