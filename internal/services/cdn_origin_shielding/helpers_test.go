package cdn_origin_shielding_test

import (
	"context"
	"fmt"
	"os"
	"strings"
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
// for use as a test fixture (origin shielding attaches to an existing CDN
// resource). The resources are cleaned up automatically via t.Cleanup.
//
// It skips the test if TF_ACC is not set, preventing failures during regular
// unit test runs.
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

	originGroup, err := client.CDN.OriginGroups.New(ctx, cdn.OriginGroupNewParams{
		OfNoneAuth: &cdn.OriginGroupNewParamsBodyNoneAuth{
			Name: name,
			Sources: []cdn.OriginGroupNewParamsBodyNoneAuthSourceUnion{
				{
					OfHostSource: &cdn.OriginGroupNewParamsBodyNoneAuthSourceHostSource{
						Source:  "example.com",
						Enabled: param.NewOpt(true),
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to create CDN origin group fixture: %s", err)
	}

	cname := fmt.Sprintf("%s.example.com", name)
	cdnResource, err := client.CDN.CDNResources.New(ctx, cdn.CDNResourceNewParams{
		Cname:       cname,
		OriginGroup: param.NewOpt(originGroup.ID),
	})
	if err != nil {
		_ = client.CDN.OriginGroups.Delete(ctx, originGroup.ID)
		t.Fatalf("failed to create CDN resource fixture: %s", err)
	}

	fixture := &cdnResourceFixture{
		ResourceID:    cdnResource.ID,
		OriginGroupID: originGroup.ID,
	}

	t.Cleanup(func() {
		if err := client.CDN.CDNResources.DeactivateAndDelete(ctx, fixture.ResourceID); err != nil {
			t.Logf("warning: failed to delete CDN resource fixture %d: %s", fixture.ResourceID, err)
		}
		if err := client.CDN.OriginGroups.Delete(ctx, fixture.OriginGroupID); err != nil {
			t.Logf("warning: failed to delete CDN origin group fixture %d: %s", fixture.OriginGroupID, err)
		}
	})

	return fixture
}

// requireOriginShielding skips the test when origin shielding is not enabled on
// the account's tariff plan (the API answers such accounts with an error that
// mentions the tariff plan). Origin shielding is a paid add-on.
func requireOriginShielding(t *testing.T, resourceID int64) {
	t.Helper()
	client, err := acctest.NewGcoreClient()
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.CDN.CDNResources.Shield.Get(context.Background(), resourceID)
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "tariff") {
		t.Skipf("origin shielding not available on this account's tariff plan: %s", err)
	}
	if err != nil {
		t.Fatalf("failed to check origin shielding availability: %s", err)
	}
}

// shieldingLocationIDs returns the origin shielding location (PoP) IDs available
// on the account, skipping the test if none exist.
func shieldingLocationIDs(t *testing.T) []int64 {
	t.Helper()
	client, err := acctest.NewGcoreClient()
	if err != nil {
		t.Fatal(err)
	}
	res, err := client.CDN.ShieldingLocation.List(context.Background(), cdn.ShieldingLocationListParams{})
	if err != nil {
		t.Fatalf("failed to list origin shielding locations: %s", err)
	}
	var ids []int64
	for _, item := range res.Results {
		ids = append(ids, item.ID)
	}
	if len(ids) == 0 {
		t.Skip("no origin shielding locations available on this account")
	}
	return ids
}
