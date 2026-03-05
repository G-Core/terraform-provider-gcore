package acctest

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/cdn"
	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/option"
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

// CDNResourceFixture holds the IDs of a CDN resource and its origin group
// created for acceptance testing. Use SetupCDNResource to create one.
type CDNResourceFixture struct {
	ResourceID    int64
	OriginGroupID int64
}

// SetupCDNResource creates a CDN origin group and CDN resource via the Go SDK
// for use as a test fixture. The resources are cleaned up automatically when
// the test completes via t.Cleanup.
//
// This is intended to be called once per test file (e.g., at the top of each
// test function) to provide a CDN resource_id that cdn_cdn_resource_rule tests
// can reference.
func SetupCDNResource(t *testing.T) *CDNResourceFixture {
	t.Helper()

	apiKey := os.Getenv("GCORE_API_KEY")
	if apiKey == "" {
		t.Fatal("GCORE_API_KEY must be set for acceptance tests")
	}

	client := gcore.NewClient(option.WithAPIKey(apiKey))
	ctx := context.Background()
	name := RandomName()

	// Create an origin group with a simple public origin
	originGroup, err := client.CDN.OriginGroups.New(ctx, cdn.OriginGroupNewParams{
		OfNoneAuth: &cdn.OriginGroupNewParamsBodyNoneAuth{
			Name: name,
			Sources: []cdn.OriginGroupNewParamsBodyNoneAuthSource{
				{
					Source:  param.NewOpt("example.com"),
					Enabled: param.NewOpt(true),
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to create CDN origin group fixture: %s", err)
	}

	// Create a CDN resource using the origin group
	cname := fmt.Sprintf("%s.example.com", name)
	cdnResource, err := client.CDN.CDNResources.New(ctx, cdn.CDNResourceNewParams{
		Cname:       cname,
		OriginGroup: param.NewOpt(originGroup.ID),
	})
	if err != nil {
		// Clean up the origin group if CDN resource creation fails
		_ = client.CDN.OriginGroups.Delete(ctx, originGroup.ID)
		t.Fatalf("failed to create CDN resource fixture: %s", err)
	}

	fixture := &CDNResourceFixture{
		ResourceID:    cdnResource.ID,
		OriginGroupID: originGroup.ID,
	}

	t.Cleanup(func() {
		// DeactivateAndDelete sets active=false then deletes
		if err := client.CDN.CDNResources.DeactivateAndDelete(ctx, fixture.ResourceID); err != nil {
			t.Logf("warning: failed to delete CDN resource fixture %d: %s", fixture.ResourceID, err)
		}
		if err := client.CDN.OriginGroups.Delete(ctx, fixture.OriginGroupID); err != nil {
			t.Logf("warning: failed to delete CDN origin group fixture %d: %s", fixture.OriginGroupID, err)
		}
	})

	return fixture
}
