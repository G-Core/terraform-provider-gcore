package cloud_load_balancer

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
	resource.AddTestSweepers("gcore_cloud_load_balancer", &resource.Sweeper{
		Name: "gcore_cloud_load_balancer",
		F:    sweepCloudLoadBalancers,
	})
}

func sweepCloudLoadBalancers(_ string) error {
	const resourceType = "gcore_cloud_load_balancer"
	defer sweep.Finish(resourceType)

	ctx := context.Background()
	target, err := sweep.Resolve(ctx)
	if err != nil {
		return err
	}

	// Collect first, delete second. Deleting while iterating an
	// offset-paginated collection shifts it, so the next offset skips exactly
	// as many items as were removed.
	var doomed []cloud.LoadBalancer
	iter := target.Client.Cloud.LoadBalancers.ListAutoPaging(ctx, cloud.LoadBalancerListParams{
		ProjectID: param.NewOpt(target.ProjectID),
		RegionID:  param.NewOpt(target.RegionID),
	})
	for iter.Next() {
		lb := iter.Current()
		sweep.RecordListed(resourceType, 1)
		if sweep.ShouldSweep(resourceType, lb.Name) {
			doomed = append(doomed, lb)
		}
	}
	if err := iter.Err(); err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping load balancer sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing load balancers: %w", err)
	}

	for _, lb := range doomed {
		log.Printf("[INFO] Deleting load balancer: %s (%s)", lb.Name, lb.ID)
		if err := target.Client.Cloud.LoadBalancers.DeleteAndPoll(ctx, lb.ID, cloud.LoadBalancerDeleteParams{
			ProjectID: param.NewOpt(target.ProjectID),
			RegionID:  param.NewOpt(target.RegionID),
		}); err != nil {
			sweep.RecordFailure(resourceType, lb.Name, err)
			continue
		}
		sweep.RecordDeleted(resourceType, lb.Name)
	}

	return nil
}
