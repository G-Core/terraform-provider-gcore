package cloud_gpu_baremetal_cluster_image

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
	resource.AddTestSweepers("gcore_cloud_gpu_baremetal_cluster_image", &resource.Sweeper{
		Name: "gcore_cloud_gpu_baremetal_cluster_image",
		F:    sweepCloudGPUBaremetalClusterImages,
	})
}

func sweepCloudGPUBaremetalClusterImages(_ string) error {
	const resourceType = "gcore_cloud_gpu_baremetal_cluster_image"
	defer sweep.Finish(resourceType)

	ctx := context.Background()
	target, err := sweep.Resolve(ctx)
	if err != nil {
		return err
	}

	// Collect first, delete second. Deleting while iterating an
	// offset-paginated collection shifts it, so the next offset skips exactly
	// as many items as were removed.
	var doomed []cloud.GPUImage
	iter := target.Client.Cloud.GPUBaremetal.Clusters.Images.ListAutoPaging(ctx, cloud.GPUBaremetalClusterImageListParams{
		ProjectID: param.NewOpt(target.ProjectID),
		RegionID:  param.NewOpt(target.RegionID),
	})
	for iter.Next() {
		image := iter.Current()
		sweep.RecordListed(resourceType, 1)
		if sweep.ShouldSweep(resourceType, image.Name) {
			doomed = append(doomed, image)
		}
	}
	if err := iter.Err(); err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping GPU baremetal cluster image sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing GPU baremetal cluster images: %w", err)
	}

	for _, image := range doomed {
		log.Printf("[INFO] Deleting GPU baremetal cluster image: %s (%s)", image.Name, image.ID)
		if err := target.Client.Cloud.GPUBaremetal.Clusters.Images.DeleteAndPoll(ctx, image.ID, cloud.GPUBaremetalClusterImageDeleteParams{
			ProjectID: param.NewOpt(target.ProjectID),
			RegionID:  param.NewOpt(target.RegionID),
		}); err != nil {
			sweep.RecordFailure(resourceType, image.Name, err)
			continue
		}
		sweep.RecordDeleted(resourceType, image.Name)
	}

	return nil
}
