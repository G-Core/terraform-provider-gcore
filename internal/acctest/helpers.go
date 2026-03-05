package acctest

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

// ProjectID returns the project ID from environment variable
// Used for constructing test configurations
func ProjectID() string {
	return os.Getenv("GCORE_CLOUD_PROJECT_ID")
}

// RegionID returns the region ID from environment variable
// Used for constructing test configurations
func RegionID() string {
	return os.Getenv("GCORE_CLOUD_REGION_ID")
}

// RandomName generates a random resource name for testing
// ensures test resources don't conflict when running concurrently
// uses RandString for short random suffix (10 alphanumeric chars)
// example: RandomName() -> "tf-test-a1b2c3d4e5"
func RandomName() string {
	return fmt.Sprintf("tf-test-%s", acctest.RandString(10))
}

var (
	latestUbuntuImageOnce sync.Once
	latestUbuntuImageID   string
	latestUbuntuImageErr  error
)

// LatestUbuntuImageID returns the ID of the latest public Ubuntu x86_64 image.
// The result is cached so the API is called at most once per test run.
// Calls t.Fatal if no suitable image is found.
func LatestUbuntuImageID(t *testing.T) string {
	t.Helper()

	latestUbuntuImageOnce.Do(func() {
		latestUbuntuImageID, latestUbuntuImageErr = fetchLatestUbuntuImageID()
	})

	if latestUbuntuImageErr != nil {
		t.Fatalf("failed to discover latest Ubuntu image: %s", latestUbuntuImageErr)
	}
	return latestUbuntuImageID
}

func fetchLatestUbuntuImageID() (string, error) {
	client, err := NewGcoreClient()
	if err != nil {
		return "", fmt.Errorf("creating API client: %w", err)
	}

	projectID, err := strconv.ParseInt(os.Getenv("GCORE_CLOUD_PROJECT_ID"), 10, 64)
	if err != nil {
		return "", fmt.Errorf("parsing GCORE_CLOUD_PROJECT_ID: %w", err)
	}
	regionID, err := strconv.ParseInt(os.Getenv("GCORE_CLOUD_REGION_ID"), 10, 64)
	if err != nil {
		return "", fmt.Errorf("parsing GCORE_CLOUD_REGION_ID: %w", err)
	}

	result, err := client.Cloud.Instances.Images.List(
		context.Background(),
		cloud.InstanceImageListParams{
			ProjectID:  param.NewOpt(projectID),
			RegionID:   param.NewOpt(regionID),
			Visibility: cloud.InstanceImageListParamsVisibilityPublic,
		},
	)
	if err != nil {
		return "", fmt.Errorf("listing images: %w", err)
	}

	var best *cloud.Image
	var bestMajor, bestMinor int

	for i := range result.Results {
		img := &result.Results[i]
		if !strings.EqualFold(img.OsDistro, "ubuntu") {
			continue
		}
		if img.Status != "active" {
			continue
		}
		if string(img.Architecture) != "x86_64" {
			continue
		}

		major, minor, err := parseVersion(img.OsVersion)
		if err != nil {
			continue
		}

		if best == nil || major > bestMajor || (major == bestMajor && minor > bestMinor) {
			best = img
			bestMajor = major
			bestMinor = minor
		}
	}

	if best == nil {
		return "", fmt.Errorf("no active public Ubuntu x86_64 image found")
	}

	return best.ID, nil
}

// parseVersion parses "24.04" into (24, 4, nil).
func parseVersion(v string) (int, int, error) {
	parts := strings.SplitN(v, ".", 2)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid version %q", v)
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}
	return major, minor, nil
}
