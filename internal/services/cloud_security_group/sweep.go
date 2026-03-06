package cloud_security_group

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
	resource.AddTestSweepers("gcore_cloud_security_group", &resource.Sweeper{
		Name: "gcore_cloud_security_group",
		F:    sweepCloudSecurityGroups,
	})
}

func sweepCloudSecurityGroups(_ string) error {
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

	page, err := client.Cloud.SecurityGroups.List(ctx, cloud.SecurityGroupListParams{
		ProjectID: param.NewOpt(projectID),
		RegionID:  param.NewOpt(regionID),
	})
	if err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping security group sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing security groups: %w", err)
	}

	for _, sg := range page.Results {
		name := sg.Name
		id := sg.ID

		if !sweep.ShouldSweep("gcore_cloud_security_group", name) {
			continue
		}

		log.Printf("[INFO] Deleting security group: %s (%s)", name, id)
		err := client.Cloud.SecurityGroups.Delete(ctx, id, cloud.SecurityGroupDeleteParams{
			ProjectID: param.NewOpt(projectID),
			RegionID:  param.NewOpt(regionID),
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete security group %s: %s", name, err)
		}
	}

	return nil
}
