package cloud_file_share

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
	resource.AddTestSweepers("gcore_cloud_file_share", &resource.Sweeper{
		Name: "gcore_cloud_file_share",
		F:    sweepCloudFileShares,
	})
}

func sweepCloudFileShares(_ string) error {
	const resourceType = "gcore_cloud_file_share"
	defer sweep.Finish(resourceType)

	ctx := context.Background()
	target, err := sweep.Resolve(ctx)
	if err != nil {
		return err
	}

	// Collect first, delete second. Deleting while iterating an
	// offset-paginated collection shifts it, so the next offset skips exactly
	// as many items as were removed.
	var doomed []cloud.FileShare
	iter := target.Client.Cloud.FileShares.ListAutoPaging(ctx, cloud.FileShareListParams{
		ProjectID: param.NewOpt(target.ProjectID),
		RegionID:  param.NewOpt(target.RegionID),
	})
	for iter.Next() {
		fileShare := iter.Current()
		sweep.RecordListed(resourceType, 1)
		if sweep.ShouldSweep(resourceType, fileShare.Name) {
			doomed = append(doomed, fileShare)
		}
	}
	if err := iter.Err(); err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping file share sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing file shares: %w", err)
	}

	for _, fileShare := range doomed {
		log.Printf("[INFO] Deleting file share: %s (%s)", fileShare.Name, fileShare.ID)
		if err := target.Client.Cloud.FileShares.DeleteAndPoll(ctx, fileShare.ID, cloud.FileShareDeleteParams{
			ProjectID: param.NewOpt(target.ProjectID),
			RegionID:  param.NewOpt(target.RegionID),
		}); err != nil {
			sweep.RecordFailure(resourceType, fileShare.Name, err)
			continue
		}
		sweep.RecordDeleted(resourceType, fileShare.Name)
	}

	return nil
}
