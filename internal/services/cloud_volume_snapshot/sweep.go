package cloud_volume_snapshot

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
	resource.AddTestSweepers("gcore_cloud_volume_snapshot", &resource.Sweeper{
		Name: "gcore_cloud_volume_snapshot",
		F:    sweepCloudVolumeSnapshots,
	})
}

func sweepCloudVolumeSnapshots(_ string) error {
	const resourceType = "gcore_cloud_volume_snapshot"
	defer sweep.Finish(resourceType)

	ctx := context.Background()
	target, err := sweep.Resolve(ctx)
	if err != nil {
		return err
	}

	// Collect first, delete second. Deleting while iterating an
	// offset-paginated collection shifts it, so the next offset skips exactly
	// as many items as were removed.
	var doomed []cloud.Snapshot
	iter := target.Client.Cloud.VolumeSnapshots.ListAutoPaging(ctx, cloud.VolumeSnapshotListParams{
		ProjectID: param.NewOpt(target.ProjectID),
		RegionID:  param.NewOpt(target.RegionID),
	})
	for iter.Next() {
		snapshot := iter.Current()
		sweep.RecordListed(resourceType, 1)
		if sweep.ShouldSweep(resourceType, snapshot.Name) {
			doomed = append(doomed, snapshot)
		}
	}
	if err := iter.Err(); err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping volume snapshot sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing volume snapshots: %w", err)
	}

	for _, snapshot := range doomed {
		log.Printf("[INFO] Deleting volume snapshot: %s (%s)", snapshot.Name, snapshot.ID)
		if err := target.Client.Cloud.VolumeSnapshots.DeleteAndPoll(ctx, snapshot.ID, cloud.VolumeSnapshotDeleteParams{
			ProjectID: param.NewOpt(target.ProjectID),
			RegionID:  param.NewOpt(target.RegionID),
		}); err != nil {
			sweep.RecordFailure(resourceType, snapshot.Name, err)
			continue
		}
		sweep.RecordDeleted(resourceType, snapshot.Name)
	}

	return nil
}
