package cloud_load_balancer

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
	resource.AddTestSweepers("gcore_cloud_load_balancer", &resource.Sweeper{
		Name: "gcore_cloud_load_balancer",
		F:    sweepCloudLoadBalancers,
	})
}

func sweepCloudLoadBalancers(_ string) error {
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

	page, err := client.Cloud.LoadBalancers.List(ctx, cloud.LoadBalancerListParams{
		ProjectID: param.NewOpt(projectID),
		RegionID:  param.NewOpt(regionID),
	})
	if err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping load balancer sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing load balancers: %w", err)
	}

	for page != nil && len(page.Results) > 0 {
		for _, lb := range page.Results {
			name := lb.Name
			id := lb.ID

			if !sweep.ShouldSweep("gcore_cloud_load_balancer", name) {
				continue
			}

			log.Printf("[INFO] Deleting load balancer: %s (%s)", name, id)
			err := client.Cloud.LoadBalancers.DeleteAndPoll(ctx, id, cloud.LoadBalancerDeleteParams{
				ProjectID: param.NewOpt(projectID),
				RegionID:  param.NewOpt(regionID),
			})
			if err != nil {
				log.Printf("[ERROR] Failed to delete load balancer %s: %s", name, err)
			}
		}

		page, err = page.GetNextPage()
		if err != nil {
			return fmt.Errorf("error getting next page of load balancers: %w", err)
		}
	}

	return nil
}
