package storage_sftp

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
	resource.AddTestSweepers("gcore_storage_sftp", &resource.Sweeper{
		Name: "gcore_storage_sftp",
		F:    sweepStorageSftp,
	})
}

func sweepStorageSftp(_ string) error {
	if err := sweep.ValidateSweeperEnvironment(); err != nil {
		return err
	}

	client := gcore.NewClient(option.WithAPIKey(os.Getenv("GCORE_API_KEY")))
	ctx := context.Background()

	page, err := client.Storage.SftpStorages.List(ctx, storage.SftpStorageListParams{})
	if err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping storage sftp sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing storage sftp: %w", err)
	}

	for _, st := range page.Results {
		if !sweep.ShouldSweep("gcore_storage_sftp", st.Name) {
			continue
		}

		log.Printf("[INFO] Deleting storage sftp: %s (%d)", st.Name, st.ID)
		if err := client.Storage.SftpStorages.DeleteAndPoll(ctx, st.ID); err != nil {
			log.Printf("[ERROR] Failed to delete storage sftp %s: %s", st.Name, err)
		}
	}

	return nil
}
