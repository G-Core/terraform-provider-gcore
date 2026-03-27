package cloud_network_router

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
	resource.AddTestSweepers("gcore_cloud_network_router", &resource.Sweeper{
		Name: "gcore_cloud_network_router",
		F:    sweepCloudNetworkRouters,
	})
}

func sweepCloudNetworkRouters(_ string) error {
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

	page, err := client.Cloud.Networks.Routers.List(ctx, cloud.NetworkRouterListParams{
		ProjectID: param.NewOpt(projectID),
		RegionID:  param.NewOpt(regionID),
	})
	if err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping router sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing routers: %w", err)
	}

	for page != nil && len(page.Results) > 0 {
		for _, router := range page.Results {
			name := router.Name
			id := router.ID

			if !sweep.ShouldSweep("gcore_cloud_network_router", name) {
				continue
			}

			log.Printf("[INFO] Deleting router: %s (%s)", name, id)
			err := client.Cloud.Networks.Routers.DeleteAndPoll(ctx, id, cloud.NetworkRouterDeleteParams{
				ProjectID: param.NewOpt(projectID),
				RegionID:  param.NewOpt(regionID),
			})
			if err != nil {
				log.Printf("[ERROR] Failed to delete router %s: %s", name, err)
			}
		}

		page, err = page.GetNextPage()
		if err != nil {
			return fmt.Errorf("error getting next page of routers: %w", err)
		}
	}

	return nil
}
