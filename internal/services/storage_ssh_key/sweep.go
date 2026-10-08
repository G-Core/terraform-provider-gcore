package storage_ssh_key

import (
	"context"
	"fmt"
	"log"

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
	const resourceType = "gcore_storage_ssh_key"
	defer sweep.Finish(resourceType)

	ctx := context.Background()
	target, err := sweep.Resolve(ctx)
	if err != nil {
		return err
	}

	// Collect first, delete second. Deleting while iterating an
	// offset-paginated collection shifts it, so the next offset skips exactly
	// as many items as were removed.
	//
	// Account-global: this enumerates every storage ssh key in the account, so
	// the token match is the only thing bounding what gets deleted.
	var doomed []storage.SSHKey
	iter := target.Client.Storage.SSHKeys.ListAutoPaging(ctx, storage.SSHKeyListParams{})
	for iter.Next() {
		key := iter.Current()
		sweep.RecordListed(resourceType, 1)
		if sweep.ShouldSweep(resourceType, key.Name) {
			doomed = append(doomed, key)
		}
	}
	if err := iter.Err(); err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping storage ssh key sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing storage ssh keys: %w", err)
	}

	for _, key := range doomed {
		log.Printf("[INFO] Deleting storage ssh key: %s (%d)", key.Name, key.ID)
		if err := target.Client.Storage.SSHKeys.Delete(ctx, key.ID); err != nil {
			sweep.RecordFailure(resourceType, key.Name, err)
			continue
		}
		sweep.RecordDeleted(resourceType, key.Name)
	}

	return nil
}
