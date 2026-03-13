package cloud_network_subnet

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/option"
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
	if err := sweep.ValidateSweeperEnvironment(); err != nil {
		return err
	}

	apiKey := os.Getenv("GCORE_API_KEY")
	projectID, err := strconv.ParseInt(os.Getenv("GCORE_CLOUD_PROJECT_ID"), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse GCORE_CLOUD_PROJECT_ID: %w", err)
	}
	regionID, err := strconv.ParseInt(os.Getenv("GCORE_CLOUD_REGION_ID"), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse GCORE_CLOUD_REGION_ID: %w", err)
	}

	client := gcore.NewClient(
		option.WithAPIKey(apiKey),
	)

	ctx := context.Background()

	page, err := client.Cloud.Networks.Subnets.List(ctx, cloud.NetworkSubnetListParams{
		ProjectID: param.NewOpt(projectID),
		RegionID:  param.NewOpt(regionID),
	})
	if err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping subnet sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing subnets: %w", err)
	}

	for page != nil && len(page.Results) > 0 {
		for _, subnet := range page.Results {
			name := subnet.Name
			id := subnet.ID

			if !sweep.ShouldSweep("gcore_cloud_network_subnet", name) {
				continue
			}

			log.Printf("[INFO] Deleting subnet: %s (%s)", name, id)
			_, err := client.Cloud.Networks.Subnets.Delete(ctx, id, cloud.NetworkSubnetDeleteParams{
				ProjectID: param.NewOpt(projectID),
				RegionID:  param.NewOpt(regionID),
			})
			if err != nil {
				log.Printf("[ERROR] Failed to delete subnet %s: %s", name, err)
			}
		}

		page, err = page.GetNextPage()
		if err != nil {
			return fmt.Errorf("error getting next page of subnets: %w", err)
		}
	}

	return nil
}
