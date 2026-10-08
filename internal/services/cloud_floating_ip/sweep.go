package cloud_floating_ip

import (
	"context"
	"fmt"
	"log"

	"github.com/G-Core/gcore-go/cloud"
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
	const resourceType = "gcore_cloud_floating_ip"
	defer sweep.Finish(resourceType)

	ctx := context.Background()
	target, err := sweep.Resolve(ctx)
	if err != nil {
		return err
	}

	// Collect first, delete second. Deleting while iterating an
	// offset-paginated collection shifts it, so the next offset skips exactly
	// as many items as were removed.
	var doomed []cloud.FloatingIPDetailed
	iter := target.Client.Cloud.FloatingIPs.ListAutoPaging(ctx, cloud.FloatingIPListParams{
		ProjectID: param.NewOpt(target.ProjectID),
		RegionID:  param.NewOpt(target.RegionID),
	})
	for iter.Next() {
		fip := iter.Current()
		sweep.RecordListed(resourceType, 1)

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

		if !sweep.ShouldSweepToken(resourceType, fip.ID, token) {
			continue
		}

		// Leave attached addresses for a later sweep, once whatever holds them
		// has itself been reclaimed.
		if fip.Status != "DOWN" {
			sweep.RecordSkipped(resourceType, fip.ID, "still attached, status "+string(fip.Status))
			continue
		}

		doomed = append(doomed, fip)
	}
	if err := iter.Err(); err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping floating IP sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing floating IPs: %w", err)
	}

	for _, fip := range doomed {
		log.Printf("[INFO] Deleting floating IP: %s", fip.ID)
		if _, err := target.Client.Cloud.FloatingIPs.Delete(ctx, fip.ID, cloud.FloatingIPDeleteParams{
			ProjectID: param.NewOpt(target.ProjectID),
			RegionID:  param.NewOpt(target.RegionID),
		}); err != nil {
			sweep.RecordFailure(resourceType, fip.ID, err)
			continue
		}
		sweep.RecordDeleted(resourceType, fip.ID)
	}

	return nil
}
