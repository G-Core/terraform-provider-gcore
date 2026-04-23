package storage_ssh_key

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/option"
	"github.com/G-Core/gcore-go/storage"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/G-Core/terraform-provider-gcore/internal/sweep"
)

func init() {
	resource.AddTestSweepers("gcore_storage_ssh_key", &resource.Sweeper{
		Name: "gcore_storage_ssh_key",
		F:    sweepStorageSSHKeys,
	})
}

func sweepStorageSSHKeys(_ string) error {
	if err := sweep.ValidateSweeperEnvironment(); err != nil {
		return err
	}

	client := gcore.NewClient(option.WithAPIKey(os.Getenv("GCORE_API_KEY")))
	ctx := context.Background()

	page, err := client.Storage.SSHKeys.List(ctx, storage.SSHKeyListParams{})
	if err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping storage ssh key sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing storage ssh keys: %w", err)
	}

	for _, key := range page.Results {
		if !sweep.ShouldSweep("gcore_storage_ssh_key", key.Name) {
			continue
		}

		log.Printf("[INFO] Deleting storage ssh key: %s (%d)", key.Name, key.ID)
		if err := client.Storage.SSHKeys.Delete(ctx, key.ID); err != nil {
			log.Printf("[ERROR] Failed to delete storage ssh key %s: %s", key.Name, err)
		}
	}

	return nil
}
