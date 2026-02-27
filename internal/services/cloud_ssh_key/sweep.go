package cloud_ssh_key

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
	resource.AddTestSweepers("gcore_cloud_ssh_key", &resource.Sweeper{
		Name: "gcore_cloud_ssh_key",
		F:    sweepCloudSSHKeys,
	})
}

func sweepCloudSSHKeys(_ string) error {
	if err := sweep.ValidateSweeperEnvironment(); err != nil {
		return err
	}

	apiKey := os.Getenv("GCORE_API_KEY")
	projectID, err := strconv.ParseInt(os.Getenv("GCORE_CLOUD_PROJECT_ID"), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse GCORE_CLOUD_PROJECT_ID: %w", err)
	}

	client := gcore.NewClient(
		option.WithAPIKey(apiKey),
		option.WithCloudProjectID(projectID),
	)

	ctx := context.Background()

	page, err := client.Cloud.SSHKeys.List(ctx, cloud.SSHKeyListParams{
		ProjectID: param.NewOpt(projectID),
	})
	if err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping SSH key sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing SSH keys: %w", err)
	}

	for _, key := range page.Results {
		name := key.Name
		id := key.ID

		if !sweep.ShouldSweep("gcore_cloud_ssh_key", name) {
			continue
		}

		log.Printf("[INFO] Deleting SSH key: %s (%s)", name, id)
		err := client.Cloud.SSHKeys.Delete(ctx, id, cloud.SSHKeyDeleteParams{
			ProjectID: param.NewOpt(projectID),
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete SSH key %s: %s", name, err)
		}
	}

	return nil
}
