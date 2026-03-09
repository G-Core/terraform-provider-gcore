package cloud_instance

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
	resource.AddTestSweepers("gcore_cloud_instance", &resource.Sweeper{
		Name:         "gcore_cloud_instance",
		F:            sweepCloudInstances,
		Dependencies: []string{"gcore_cloud_volume"},
	})
}

func sweepCloudInstances(_ string) error {
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

	instances, err := client.Cloud.Instances.List(ctx, cloud.InstanceListParams{
		ProjectID: param.NewOpt(projectID),
		RegionID:  param.NewOpt(regionID),
	})
	if err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping cloud instance sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing cloud instances: %w", err)
	}

	for _, instance := range instances.Results {
		if !sweep.ShouldSweep("gcore_cloud_instance", instance.Name) {
			continue
		}

		log.Printf("[INFO] Deleting cloud instance: %s (%s)", instance.Name, instance.ID)
		err := client.Cloud.Instances.DeleteAndPoll(ctx, instance.ID, cloud.InstanceDeleteParams{
			ProjectID: param.NewOpt(projectID),
			RegionID:  param.NewOpt(regionID),
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete cloud instance %s: %s", instance.Name, err)
		}
	}

	return nil
}
