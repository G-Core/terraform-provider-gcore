package cdn_origin_group

import (
	"context"
	"fmt"
	"log"

	"github.com/G-Core/gcore-go/cdn"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/G-Core/terraform-provider-gcore/internal/sweep"
)

func init() {
	resource.AddTestSweepers("gcore_cdn_origin_group", &resource.Sweeper{
		Name: "gcore_cdn_origin_group",
		F:    sweepCDNOriginGroups,
	})
}

func sweepCDNOriginGroups(_ string) error {
	const resourceType = "gcore_cdn_origin_group"
	defer sweep.Finish(resourceType)

	ctx := context.Background()
	target, err := sweep.Resolve(ctx)
	if err != nil {
		return err
	}

	// Collect first, delete second. Deleting while iterating an
	// offset-paginated collection shifts it, so the next offset skips exactly
	// as many items as were removed.
	//
	// Account-global: this enumerates every origin group in the account, so the
	// token match is the only thing bounding what gets deleted.
	var doomed []cdn.OriginGroupsUnion
	iter := target.Client.CDN.OriginGroups.ListAutoPaging(ctx, cdn.OriginGroupListParams{})
	for iter.Next() {
		group := iter.Current()
		sweep.RecordListed(resourceType, 1)
		if !sweep.ShouldSweep(resourceType, group.Name) {
			continue
		}

		// Skip groups with related resources - cannot delete if CDN resources are using it
		if group.HasRelatedResources {
			sweep.RecordSkipped(resourceType, group.Name, "still used by a CDN resource")
			continue
		}

		doomed = append(doomed, group)
	}
	if err := iter.Err(); err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping CDN origin group sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing CDN origin groups: %w", err)
	}

	for _, group := range doomed {
		log.Printf("[INFO] Deleting CDN origin group: %s (%d)", group.Name, group.ID)
		if err := target.Client.CDN.OriginGroups.Delete(ctx, group.ID); err != nil {
			sweep.RecordFailure(resourceType, group.Name, err)
			continue
		}
		sweep.RecordDeleted(resourceType, group.Name)
	}

	return nil
}
