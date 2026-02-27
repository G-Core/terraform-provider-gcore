package cloud_network

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

	"github.com/stainless-sdks/gcore-terraform/internal/sweep"
)

func init() {
	resource.AddTestSweepers("gcore_cloud_network", &resource.Sweeper{
		Name: "gcore_cloud_network",
		F:    sweepCloudNetworks,
	})
}

func sweepCloudNetworks(_ string) error {
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

	page, err := client.Cloud.Networks.List(ctx, cloud.NetworkListParams{
		ProjectID: param.NewOpt(projectID),
		RegionID:  param.NewOpt(regionID),
	})
	if err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping network sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing networks: %w", err)
	}

	for _, network := range page.Results {
		name := network.Name
		id := network.ID

		if !sweep.ShouldSweep("gcore_cloud_network", name) {
			continue
		}

		log.Printf("[INFO] Deleting network: %s (%s)", name, id)
		_, err := client.Cloud.Networks.Delete(ctx, id, cloud.NetworkDeleteParams{
			ProjectID: param.NewOpt(projectID),
			RegionID:  param.NewOpt(regionID),
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete network %s: %s", name, err)
		}
	}

	return nil
}
