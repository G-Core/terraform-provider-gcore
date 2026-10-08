package cloud_instance

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
	resource.AddTestSweepers("gcore_cloud_instance", &resource.Sweeper{
		Name:         "gcore_cloud_instance",
		F:            sweepCloudInstances,
		Dependencies: []string{"gcore_cloud_volume"},
	})
}

func sweepCloudInstances(_ string) error {
	const resourceType = "gcore_cloud_instance"
	defer sweep.Finish(resourceType)

	ctx := context.Background()
	target, err := sweep.Resolve(ctx)
	if err != nil {
		return err
	}

	// Collect first, delete second. Deleting while iterating an
	// offset-paginated collection shifts it, so the next offset skips exactly
	// as many items as were removed.
	var doomed []cloud.Instance
	iter := target.Client.Cloud.Instances.ListAutoPaging(ctx, cloud.InstanceListParams{
		ProjectID: param.NewOpt(target.ProjectID),
		RegionID:  param.NewOpt(target.RegionID),
	})
	for iter.Next() {
		instance := iter.Current()
		sweep.RecordListed(resourceType, 1)
		if sweep.ShouldSweep(resourceType, instance.Name) {
			doomed = append(doomed, instance)
		}
	}
	if err := iter.Err(); err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping cloud instance sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing cloud instances: %w", err)
	}

	for _, instance := range doomed {
		log.Printf("[INFO] Deleting cloud instance: %s (%s)", instance.Name, instance.ID)
		if err := target.Client.Cloud.Instances.DeleteAndPoll(ctx, instance.ID, cloud.InstanceDeleteParams{
			ProjectID: param.NewOpt(target.ProjectID),
			RegionID:  param.NewOpt(target.RegionID),
		}); err != nil {
			sweep.RecordFailure(resourceType, instance.Name, err)
			continue
		}
		sweep.RecordDeleted(resourceType, instance.Name)
	}

	return nil
}
