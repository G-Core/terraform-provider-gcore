package cloud_volume

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
	resource.AddTestSweepers("gcore_cloud_volume", &resource.Sweeper{
		Name: "gcore_cloud_volume",
		F:    sweepCloudVolumes,
		// Snapshots must be swept first: a volume cannot be deleted while
		// snapshots created from it still exist.
		Dependencies: []string{"gcore_cloud_volume_snapshot"},
	})
}

func sweepCloudVolumes(_ string) error {
	const resourceType = "gcore_cloud_volume"
	defer sweep.Finish(resourceType)

	ctx := context.Background()
	target, err := sweep.Resolve(ctx)
	if err != nil {
		return err
	}

	// Collect first, delete second. Deleting while iterating an
	// offset-paginated collection shifts it, so the next offset skips exactly
	// as many items as were removed.
	var doomed []cloud.Volume
	iter := target.Client.Cloud.Volumes.ListAutoPaging(ctx, cloud.VolumeListParams{
		ProjectID: param.NewOpt(target.ProjectID),
		RegionID:  param.NewOpt(target.RegionID),
	})
	for iter.Next() {
		volume := iter.Current()
		sweep.RecordListed(resourceType, 1)
		if sweep.ShouldSweep(resourceType, volume.Name) {
			doomed = append(doomed, volume)
		}
	}
	if err := iter.Err(); err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping volume sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing volumes: %w", err)
	}

	for _, volume := range doomed {
		log.Printf("[INFO] Deleting volume: %s (%s)", volume.Name, volume.ID)
		if err := target.Client.Cloud.Volumes.DeleteAndPoll(ctx, volume.ID, cloud.VolumeDeleteParams{
			ProjectID: param.NewOpt(target.ProjectID),
			RegionID:  param.NewOpt(target.RegionID),
		}); err != nil {
			sweep.RecordFailure(resourceType, volume.Name, err)
			continue
		}
		sweep.RecordDeleted(resourceType, volume.Name)
	}

	return nil
}
