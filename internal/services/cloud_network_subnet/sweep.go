package cloud_network_subnet

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
	resource.AddTestSweepers("gcore_cloud_network_subnet", &resource.Sweeper{
		Name:         "gcore_cloud_network_subnet",
		F:            sweepCloudNetworkSubnets,
		Dependencies: []string{"gcore_cloud_instance"},
	})
}

func sweepCloudNetworkSubnets(_ string) error {
	const resourceType = "gcore_cloud_network_subnet"
	defer sweep.Finish(resourceType)

	ctx := context.Background()
	target, err := sweep.Resolve(ctx)
	if err != nil {
		return err
	}

	// Collect first, delete second. Deleting while iterating an
	// offset-paginated collection shifts it, so the next offset skips exactly
	// as many items as were removed.
	var doomed []cloud.Subnet
	iter := target.Client.Cloud.Networks.Subnets.ListAutoPaging(ctx, cloud.NetworkSubnetListParams{
		ProjectID: param.NewOpt(target.ProjectID),
		RegionID:  param.NewOpt(target.RegionID),
	})
	for iter.Next() {
		subnet := iter.Current()
		sweep.RecordListed(resourceType, 1)
		if sweep.ShouldSweep(resourceType, subnet.Name) {
			doomed = append(doomed, subnet)
		}
	}
	if err := iter.Err(); err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping subnet sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing subnets: %w", err)
	}

	for _, subnet := range doomed {
		log.Printf("[INFO] Deleting subnet: %s (%s)", subnet.Name, subnet.ID)
		if _, err := target.Client.Cloud.Networks.Subnets.Delete(ctx, subnet.ID, cloud.NetworkSubnetDeleteParams{
			ProjectID: param.NewOpt(target.ProjectID),
			RegionID:  param.NewOpt(target.RegionID),
		}); err != nil {
			sweep.RecordFailure(resourceType, subnet.Name, err)
			continue
		}
		sweep.RecordDeleted(resourceType, subnet.Name)
	}

	return nil
}
