package storage_sftp

import (
	"context"
	"fmt"
	"log"

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
	const resourceType = "gcore_storage_sftp"
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
	// Account-global: this enumerates every sftp storage in the account, so the
	// token match is the only thing bounding what gets deleted.
	var doomed []storage.SftpStorage
	iter := target.Client.Storage.SftpStorages.ListAutoPaging(ctx, storage.SftpStorageListParams{})
	for iter.Next() {
		st := iter.Current()
		sweep.RecordListed(resourceType, 1)
		if sweep.ShouldSweep(resourceType, st.Name) {
			doomed = append(doomed, st)
		}
	}
	if err := iter.Err(); err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping storage sftp sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing storage sftp: %w", err)
	}

	for _, st := range doomed {
		log.Printf("[INFO] Deleting storage sftp: %s (%d)", st.Name, st.ID)
		if err := target.Client.Storage.SftpStorages.DeleteAndPoll(ctx, st.ID); err != nil {
			sweep.RecordFailure(resourceType, st.Name, err)
			continue
		}
		sweep.RecordDeleted(resourceType, st.Name)
	}

	return nil
}
