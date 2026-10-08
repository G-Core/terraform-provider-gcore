package cloud_network_router

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
	resource.AddTestSweepers("gcore_cloud_network_router", &resource.Sweeper{
		Name: "gcore_cloud_network_router",
		F:    sweepCloudNetworkRouters,
	})
}

func sweepCloudNetworkRouters(_ string) error {
	const resourceType = "gcore_cloud_network_router"
	defer sweep.Finish(resourceType)

	ctx := context.Background()
	target, err := sweep.Resolve(ctx)
	if err != nil {
		return err
	}

	// Collect first, delete second. Deleting while iterating an
	// offset-paginated collection shifts it, so the next offset skips exactly
	// as many items as were removed.
	var doomed []cloud.Router
	iter := target.Client.Cloud.Networks.Routers.ListAutoPaging(ctx, cloud.NetworkRouterListParams{
		ProjectID: param.NewOpt(target.ProjectID),
		RegionID:  param.NewOpt(target.RegionID),
	})
	for iter.Next() {
		router := iter.Current()
		sweep.RecordListed(resourceType, 1)
		if sweep.ShouldSweep(resourceType, router.Name) {
			doomed = append(doomed, router)
		}
	}
	if err := iter.Err(); err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping router sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing routers: %w", err)
	}

	for _, router := range doomed {
		log.Printf("[INFO] Deleting router: %s (%s)", router.Name, router.ID)
		if err := target.Client.Cloud.Networks.Routers.DeleteAndPoll(ctx, router.ID, cloud.NetworkRouterDeleteParams{
			ProjectID: param.NewOpt(target.ProjectID),
			RegionID:  param.NewOpt(target.RegionID),
		}); err != nil {
			sweep.RecordFailure(resourceType, router.Name, err)
			continue
		}
		sweep.RecordDeleted(resourceType, router.Name)
	}

	return nil
}
