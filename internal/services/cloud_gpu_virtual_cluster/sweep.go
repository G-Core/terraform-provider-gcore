package cloud_gpu_virtual_cluster

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
	resource.AddTestSweepers("gcore_cloud_gpu_virtual_cluster", &resource.Sweeper{
		Name: "gcore_cloud_gpu_virtual_cluster",
		F:    sweepCloudGPUVirtualClusters,
	})
}

func sweepCloudGPUVirtualClusters(_ string) error {
	const resourceType = "gcore_cloud_gpu_virtual_cluster"
	defer sweep.Finish(resourceType)

	ctx := context.Background()
	target, err := sweep.Resolve(ctx)
	if err != nil {
		return err
	}

	// Collect first, delete second. Deleting while iterating an
	// offset-paginated collection shifts it, so the next offset skips exactly
	// as many items as were removed.
	var doomed []cloud.GPUVirtualCluster
	iter := target.Client.Cloud.GPUVirtual.Clusters.ListAutoPaging(ctx, cloud.GPUVirtualClusterListParams{
		ProjectID: param.NewOpt(target.ProjectID),
		RegionID:  param.NewOpt(target.RegionID),
	})
	for iter.Next() {
		cluster := iter.Current()
		sweep.RecordListed(resourceType, 1)
		if sweep.ShouldSweep(resourceType, cluster.Name) {
			doomed = append(doomed, cluster)
		}
	}
	if err := iter.Err(); err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping GPU virtual cluster sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing GPU virtual clusters: %w", err)
	}

	for _, cluster := range doomed {
		log.Printf("[INFO] Deleting GPU virtual cluster: %s (%s)", cluster.Name, cluster.ID)
		if err := target.Client.Cloud.GPUVirtual.Clusters.DeleteAndPoll(ctx, cluster.ID, cloud.GPUVirtualClusterDeleteParams{
			ProjectID: param.NewOpt(target.ProjectID),
			RegionID:  param.NewOpt(target.RegionID),
		}); err != nil {
			sweep.RecordFailure(resourceType, cluster.Name, err)
			continue
		}
		sweep.RecordDeleted(resourceType, cluster.Name)
	}

	return nil
}
