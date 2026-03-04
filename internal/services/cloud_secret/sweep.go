package cloud_secret

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
	resource.AddTestSweepers("gcore_cloud_secret", &resource.Sweeper{
		Name: "gcore_cloud_secret",
		F:    sweepCloudSecrets,
	})
}

func sweepCloudSecrets(_ string) error {
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

	page, err := client.Cloud.Secrets.List(ctx, cloud.SecretListParams{
		ProjectID: param.NewOpt(projectID),
		RegionID:  param.NewOpt(regionID),
	})
	if err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping secret sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing secrets: %w", err)
	}

	for _, secret := range page.Results {
		name := secret.Name
		id := secret.ID

		if !sweep.ShouldSweep("gcore_cloud_secret", name) {
			continue
		}

		log.Printf("[INFO] Deleting secret: %s (%s)", name, id)
		err := client.Cloud.Secrets.DeleteAndPoll(ctx, id, cloud.SecretDeleteParams{
			ProjectID: param.NewOpt(projectID),
			RegionID:  param.NewOpt(regionID),
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete secret %s: %s", name, err)
		}
	}

	return nil
}
