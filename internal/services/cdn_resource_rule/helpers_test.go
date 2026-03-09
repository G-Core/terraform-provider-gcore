package cdn_resource_rule_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/cdn"
	"github.com/G-Core/gcore-go/option"
	"github.com/G-Core/gcore-go/packages/param"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

// cdnResourceFixture holds the IDs of a CDN resource and its origin group
// created for acceptance testing.
type cdnResourceFixture struct {
	ResourceID    int64
	OriginGroupID int64
}

// setupCDNResource creates a CDN origin group and CDN resource via the Go SDK
// for use as a test fixture. The resources are cleaned up automatically when
// the test completes via t.Cleanup.
//
// This is intended to be called once per test function to provide a CDN
// resource_id that cdn_cdn_resource_rule tests can reference.
//
// It skips the test if TF_ACC is not set, preventing failures during
// regular unit test runs.
func setupCDNResource(t *testing.T) *cdnResourceFixture {
	t.Helper()

	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	apiKey := os.Getenv("GCORE_API_KEY")
	if apiKey == "" {
		t.Fatal("GCORE_API_KEY must be set for acceptance tests")
	}

	client := gcore.NewClient(option.WithAPIKey(apiKey))
	ctx := context.Background()
	name := acctest.RandomName()

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

	fixture := &cdnResourceFixture{
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
