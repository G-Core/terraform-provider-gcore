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
)

var (
	latestUbuntuImageOnce sync.Once
	cachedUbuntuImage     cloud.Image
	cachedUbuntuImageErr  error
)

// LatestUbuntuImage returns the latest public Ubuntu x86_64 image available in
// the configured project and region. The result is cached so the API is called
// at most once per test binary.
//
// It skips the test if TF_ACC is not set, preventing failures during regular
// unit test runs.
func LatestUbuntuImage(t *testing.T) cloud.Image {
	t.Helper()

	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	latestUbuntuImageOnce.Do(func() {
		cachedUbuntuImage, cachedUbuntuImageErr = fetchLatestUbuntuImage()
	})

	if cachedUbuntuImageErr != nil {
		t.Fatalf("failed to discover latest Ubuntu image: %s", cachedUbuntuImageErr)
	}
	return cachedUbuntuImage
}

// LatestUbuntuImageID returns the ID of the image found by LatestUbuntuImage.
func LatestUbuntuImageID(t *testing.T) string {
	t.Helper()
	return LatestUbuntuImage(t).ID
}

func fetchLatestUbuntuImage() (cloud.Image, error) {
	var zero cloud.Image

	client, err := NewGcoreClient()
	if err != nil {
		return zero, fmt.Errorf("creating API client: %w", err)
	}

	projectID, err := strconv.ParseInt(os.Getenv("GCORE_CLOUD_PROJECT_ID"), 10, 64)
	if err != nil {
		return zero, fmt.Errorf("parsing GCORE_CLOUD_PROJECT_ID: %w", err)
	}
	regionID, err := strconv.ParseInt(os.Getenv("GCORE_CLOUD_REGION_ID"), 10, 64)
	if err != nil {
		return zero, fmt.Errorf("parsing GCORE_CLOUD_REGION_ID: %w", err)
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
		return zero, fmt.Errorf("listing images: %w", err)
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

		major, minor, err := parseImageVersion(img.OsVersion)
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
		return zero, fmt.Errorf("no active public Ubuntu x86_64 image found")
	}

	return *best, nil
}

// parseImageVersion parses "24.04" into (24, 4, nil).
func parseImageVersion(v string) (int, int, error) {
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
