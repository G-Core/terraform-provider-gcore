package cloud_network

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
	resource.AddTestSweepers("gcore_cloud_network", &resource.Sweeper{
		Name:         "gcore_cloud_network",
		F:            sweepCloudNetworks,
		Dependencies: []string{"gcore_cloud_network_subnet"},
	})
}

func sweepCloudNetworks(_ string) error {
	const resourceType = "gcore_cloud_network"
	defer sweep.Finish(resourceType)

	ctx := context.Background()
	target, err := sweep.Resolve(ctx)
	if err != nil {
		return err
	}

	// Collect first, delete second. Deleting while iterating an
	// offset-paginated collection shifts it, so the next offset skips exactly
	// as many items as were removed.
	var doomed []cloud.Network
	iter := target.Client.Cloud.Networks.ListAutoPaging(ctx, cloud.NetworkListParams{
		ProjectID: param.NewOpt(target.ProjectID),
		RegionID:  param.NewOpt(target.RegionID),
	})
	for iter.Next() {
		network := iter.Current()
		sweep.RecordListed(resourceType, 1)
		if sweep.ShouldSweep(resourceType, network.Name) {
			doomed = append(doomed, network)
		}
	}
	if err := iter.Err(); err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping network sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing networks: %w", err)
	}

	for _, network := range doomed {
		log.Printf("[INFO] Deleting network: %s (%s)", network.Name, network.ID)
		if _, err := target.Client.Cloud.Networks.Delete(ctx, network.ID, cloud.NetworkDeleteParams{
			ProjectID: param.NewOpt(target.ProjectID),
			RegionID:  param.NewOpt(target.RegionID),
		}); err != nil {
			sweep.RecordFailure(resourceType, network.Name, err)
			continue
		}
		sweep.RecordDeleted(resourceType, network.Name)
	}

	return nil
}
