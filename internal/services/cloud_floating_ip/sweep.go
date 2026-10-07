package cloud_floating_ip

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/option"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/G-Core/terraform-provider-gcore/internal/sweep"
)

func init() {
	resource.AddTestSweepers("gcore_cloud_floating_ip", &resource.Sweeper{
		Name: "gcore_cloud_floating_ip",
		F:    sweepCloudFloatingIPs,
	})
}

func sweepCloudFloatingIPs(_ string) error {
	if err := sweep.ValidateSweeperEnvironment(); err != nil {
		return err
	}

	apiKey := os.Getenv("GCORE_API_KEY")
	projectID, err := strconv.ParseInt(os.Getenv("GCORE_CLOUD_PROJECT_ID"), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse GCORE_CLOUD_PROJECT_ID: %w", err)
	}
	regionID, err := strconv.ParseInt(os.Getenv("GCORE_CLOUD_REGION_ID"), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse GCORE_CLOUD_REGION_ID: %w", err)
	}

	client := gcore.NewClient(
		option.WithAPIKey(apiKey),
	)

	ctx := context.Background()

	page, err := client.Cloud.FloatingIPs.List(ctx, cloud.FloatingIPListParams{
		ProjectID: param.NewOpt(projectID),
		RegionID:  param.NewOpt(regionID),
	})
	if err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping floating IP sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing floating IPs: %w", err)
	}

	for _, fip := range page.Results {
		id := fip.ID

		// Floating IPs have no name field, so the run token travels in a tag
		// the tests set. Nothing else distinguishes a test address from a
		// reserved production one: an earlier version deleted every unattached
		// address it found, and unattached is the normal state for an address
		// that has been reserved but not yet assigned.
		var token string
		for _, tag := range fip.Tags {
			if tag.Key == sweep.TagKey {
				token = tag.Value
				break
			}
		}

		if !sweep.ShouldSweepToken("gcore_cloud_floating_ip", id, token) {
			continue
		}

		// Leave attached addresses for a later sweep, once whatever holds them
		// has itself been reclaimed.
		if fip.Status != "DOWN" {
			log.Printf("[INFO] Skipping floating IP %s: still attached (status: %s)", id, fip.Status)
			continue
		}

		log.Printf("[INFO] Deleting floating IP: %s", id)
		_, err := client.Cloud.FloatingIPs.Delete(ctx, id, cloud.FloatingIPDeleteParams{
			ProjectID: param.NewOpt(projectID),
			RegionID:  param.NewOpt(regionID),
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete floating IP %s: %s", id, err)
		}
	}

	return nil
}
