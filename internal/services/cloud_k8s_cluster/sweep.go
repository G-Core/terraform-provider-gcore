package cloud_k8s_cluster

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
	resource.AddTestSweepers("gcore_cloud_k8s_cluster", &resource.Sweeper{
		Name: "gcore_cloud_k8s_cluster",
		F:    sweepCloudK8SClusters,
	})
}

func sweepCloudK8SClusters(_ string) error {
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

	list, err := client.Cloud.K8S.Clusters.List(ctx, cloud.K8SClusterListParams{
		ProjectID: param.NewOpt(projectID),
		RegionID:  param.NewOpt(regionID),
	})
	if err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping k8s cluster sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing k8s clusters: %w", err)
	}

	for _, cluster := range list.Results {
		name := cluster.Name

		if !sweep.ShouldSweep("gcore_cloud_k8s_cluster", name) {
			continue
		}

		log.Printf("[INFO] Deleting k8s cluster: %s", name)
		err := client.Cloud.K8S.Clusters.DeleteAndPoll(ctx, name, cloud.K8SClusterDeleteParams{
			ProjectID: param.NewOpt(projectID),
			RegionID:  param.NewOpt(regionID),
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete k8s cluster %s: %s", name, err)
		}
	}

	return nil
}
