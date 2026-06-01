package cloud_instance_image

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
	resource.AddTestSweepers("gcore_cloud_instance_image", &resource.Sweeper{
		Name: "gcore_cloud_instance_image",
		F:    sweepCloudInstanceImages,
	})
}

func sweepCloudInstanceImages(_ string) error {
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

	page, err := client.Cloud.Instances.Images.List(ctx, cloud.InstanceImageListParams{
		ProjectID: param.NewOpt(projectID),
		RegionID:  param.NewOpt(regionID),
	})
	if err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping instance image sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing instance images: %w", err)
	}

	for _, image := range page.Results {
		name := image.Name
		id := image.ID

		if !sweep.ShouldSweep("gcore_cloud_instance_image", name) {
			continue
		}

		log.Printf("[INFO] Deleting instance image: %s (%s)", name, id)
		err := client.Cloud.Instances.Images.DeleteAndPoll(ctx, id, cloud.InstanceImageDeleteParams{
			ProjectID: param.NewOpt(projectID),
			RegionID:  param.NewOpt(regionID),
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete instance image %s: %s", name, err)
		}
	}

	return nil
}
