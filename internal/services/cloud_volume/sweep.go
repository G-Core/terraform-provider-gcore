package cloud_volume

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
	resource.AddTestSweepers("gcore_cloud_volume", &resource.Sweeper{
		Name: "gcore_cloud_volume",
		F:    sweepCloudVolumes,
	})
}

func sweepCloudVolumes(_ string) error {
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
		option.WithCloudProjectID(projectID),
	)

	ctx := context.Background()

	page, err := client.Cloud.Volumes.List(ctx, cloud.VolumeListParams{
		ProjectID: param.NewOpt(projectID),
		RegionID:  param.NewOpt(regionID),
	})
	if err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping volume sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing volumes: %w", err)
	}

	for _, volume := range page.Results {
		name := volume.Name
		id := volume.ID

		if !sweep.ShouldSweep("gcore_cloud_volume", name) {
			continue
		}

		log.Printf("[INFO] Deleting volume: %s (%s)", name, id)
		err := client.Cloud.Volumes.DeleteAndPoll(ctx, id, cloud.VolumeDeleteParams{
			ProjectID: param.NewOpt(projectID),
			RegionID:  param.NewOpt(regionID),
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete volume %s: %s", name, err)
		}
	}

	return nil
}
