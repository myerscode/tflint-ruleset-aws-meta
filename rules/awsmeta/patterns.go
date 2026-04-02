package awsmeta

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/myerscode/aws-meta/pkg/partitions"
	"github.com/myerscode/aws-meta/pkg/regions"
)

var (
	regionPattern     *regexp.Regexp
	regionPatternOnce sync.Once

	regionInStringPattern     *regexp.Regexp
	regionInStringPatternOnce sync.Once

	availabilityZonePattern     *regexp.Regexp
	availabilityZonePatternOnce sync.Once

	arnRegionPattern     *regexp.Regexp
	arnRegionPatternOnce sync.Once

	partitionPattern     *regexp.Regexp
	partitionPatternOnce sync.Once

	dnsSuffixPattern     *regexp.Regexp
	dnsSuffixPatternOnce sync.Once

	accountIDPattern     *regexp.Regexp
	accountIDPatternOnce sync.Once

	amiIDPattern     *regexp.Regexp
	amiIDPatternOnce sync.Once
)

func loadRegionNames() []string {
	regionList, err := regions.ListAllRegions()
	if err != nil {
		panic(fmt.Sprintf("failed to load AWS regions: %v", err))
	}

	if len(regionList) == 0 {
		panic("AWS region list is empty")
	}

	regionNames := make([]string, 0, len(regionList))
	for _, region := range regionList {
		regionNames = append(regionNames, regexp.QuoteMeta(region.RegionId))
	}
	return regionNames
}

// GetRegionPattern returns a compiled regex pattern matching all AWS regions
func GetRegionPattern() *regexp.Regexp {
	regionPatternOnce.Do(func() {
		regionNames := loadRegionNames()
		pattern := fmt.Sprintf("^(%s)$", strings.Join(regionNames, "|"))
		regionPattern = regexp.MustCompile(pattern)
	})
	return regionPattern
}

// GetRegionInStringPattern returns a compiled regex pattern for finding regions within strings
func GetRegionInStringPattern() *regexp.Regexp {
	regionInStringPatternOnce.Do(func() {
		regionNames := loadRegionNames()
		pattern := fmt.Sprintf("(%s)", strings.Join(regionNames, "|"))
		regionInStringPattern = regexp.MustCompile(pattern)
	})
	return regionInStringPattern
}

// GetAvailabilityZonePattern returns a compiled regex pattern matching availability zones
func GetAvailabilityZonePattern() *regexp.Regexp {
	availabilityZonePatternOnce.Do(func() {
		regionNames := loadRegionNames()
		pattern := fmt.Sprintf("^(%s)[a-z]$", strings.Join(regionNames, "|"))
		availabilityZonePattern = regexp.MustCompile(pattern)
	})
	return availabilityZonePattern
}

// GetARNRegionPattern returns a compiled regex pattern for finding regions in ARNs
func GetARNRegionPattern() *regexp.Regexp {
	arnRegionPatternOnce.Do(func() {
		regionNames := loadRegionNames()
		pattern := fmt.Sprintf(`arn:aws[^:]*:[^:]+:(%s):`, strings.Join(regionNames, "|"))
		arnRegionPattern = regexp.MustCompile(pattern)
	})
	return arnRegionPattern
}

// GetPartitionPattern returns a compiled regex pattern matching all AWS partitions
func GetPartitionPattern() *regexp.Regexp {
	partitionPatternOnce.Do(func() {
		partitionList, err := partitions.List()
		if err != nil {
			panic(fmt.Sprintf("failed to load AWS partitions: %v", err))
		}

		if len(partitionList) == 0 {
			panic("AWS partition list is empty")
		}

		partitionNames := make([]string, 0, len(partitionList))
		for _, partition := range partitionList {
			partitionNames = append(partitionNames, regexp.QuoteMeta(partition.ID))
		}

		pattern := fmt.Sprintf(`arn:(%s):`, strings.Join(partitionNames, "|"))
		partitionPattern = regexp.MustCompile(pattern)
	})
	return partitionPattern
}

// GetDNSSuffixPattern returns a compiled regex pattern matching service principals
// with any known AWS DNS suffix (e.g. "s3.amazonaws.com", "lambda.amazonaws.com.cn")
func GetDNSSuffixPattern() *regexp.Regexp {
	dnsSuffixPatternOnce.Do(func() {
		partitionList, err := partitions.List()
		if err != nil {
			panic(fmt.Sprintf("failed to load AWS partitions: %v", err))
		}

		if len(partitionList) == 0 {
			panic("AWS partition list is empty")
		}

		// Collect unique DNS suffixes from all partitions
		seen := make(map[string]bool)
		var suffixes []string
		for _, partition := range partitionList {
			if partition.DNSSuffix != "" && !seen[partition.DNSSuffix] {
				seen[partition.DNSSuffix] = true
				suffixes = append(suffixes, regexp.QuoteMeta(partition.DNSSuffix))
			}
		}

		// Match: service-name.dns-suffix (e.g. s3.amazonaws.com, lambda.c2s.ic.gov)
		pattern := fmt.Sprintf(`([a-z0-9\-]+)\.(%s)`, strings.Join(suffixes, "|"))
		dnsSuffixPattern = regexp.MustCompile(pattern)
	})
	return dnsSuffixPattern
}

// GetAccountIDPattern returns a compiled regex pattern matching 12-digit AWS account IDs.
func GetAccountIDPattern() *regexp.Regexp {
	accountIDPatternOnce.Do(func() {
		// AWS account IDs are always exactly 12 digits
		accountIDPattern = regexp.MustCompile(`\b(\d{12})\b`)
	})
	return accountIDPattern
}

// GetAMIIDPattern returns a compiled regex pattern matching AWS AMI IDs.
func GetAMIIDPattern() *regexp.Regexp {
	amiIDPatternOnce.Do(func() {
		// AMI IDs follow the format ami-<hex string>
		amiIDPattern = regexp.MustCompile(`\bami-[0-9a-f]{8,17}\b`)
	})
	return amiIDPattern
}
