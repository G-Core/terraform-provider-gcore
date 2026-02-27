package cloud_placement_group

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
	resource.AddTestSweepers("gcore_cloud_placement_group", &resource.Sweeper{
		Name: "gcore_cloud_placement_group",
		F:    sweepCloudPlacementGroups,
	})
}

func sweepCloudPlacementGroups(_ string) error {
	if err := sweep.ValidateSweeperEnvironment(); err != nil {
		return err
	}

	apiKey := os.Getenv("GCORE_API_KEY")
	projectIDStr := os.Getenv("GCORE_CLOUD_PROJECT_ID")
	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse GCORE_CLOUD_PROJECT_ID: %w", err)
	}
	regionIDStr := os.Getenv("GCORE_CLOUD_REGION_ID")
	regionID, err := strconv.ParseInt(regionIDStr, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse GCORE_CLOUD_REGION_ID: %w", err)
	}

	client := gcore.NewClient(
		option.WithAPIKey(apiKey),
		option.WithCloudProjectID(projectID),
	)

	ctx := context.Background()

	groups, err := client.Cloud.PlacementGroups.List(ctx, cloud.PlacementGroupListParams{
		ProjectID: param.NewOpt(projectID),
		RegionID:  param.NewOpt(regionID),
	})
	if err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping placement group sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing placement groups: %w", err)
	}

	for _, group := range groups.Results {
		name := group.Name
		id := group.ServergroupID

		if !sweep.ShouldSweep("gcore_cloud_placement_group", name) {
			continue
		}

		log.Printf("[INFO] Deleting placement group: %s (%s)", name, id)
		_, err := client.Cloud.PlacementGroups.Delete(ctx, id, cloud.PlacementGroupDeleteParams{
			ProjectID: param.NewOpt(projectID),
			RegionID:  param.NewOpt(regionID),
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete placement group %s: %s", name, err)
		}
	}

	return nil
}
