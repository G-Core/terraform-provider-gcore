package storage_object_storage

import (
	"context"
	"fmt"
	"log"

	"github.com/G-Core/gcore-go/storage"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/G-Core/terraform-provider-gcore/internal/sweep"
)

func init() {
	resource.AddTestSweepers("gcore_storage_object_storage", &resource.Sweeper{
		Name: "gcore_storage_object_storage",
		F:    sweepStorageObjectStorage,
		// A storage bound to a CDN origin or logs target cannot be deleted (409)
		Dependencies: []string{"gcore_cdn_origin_group", "gcore_cdn_logs_uploader_target"},
	})
}

// sweepStorageObjectStorage deletes leaked object storages. Deleting a storage also
// deletes its buckets, so this sweeper reclaims gcore_storage_object_storage_bucket too.
func sweepStorageObjectStorage(_ string) error {
	const resourceType = "gcore_storage_object_storage"
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
	// Account-global: this enumerates every object storage in the account, so the
	// token match is the only thing bounding what gets deleted.
	var doomed []storage.S3Storage
	iter := target.Client.Storage.ObjectStorages.ListAutoPaging(ctx, storage.ObjectStorageListParams{})
	for iter.Next() {
		st := iter.Current()
		sweep.RecordListed(resourceType, 1)
		// A storage already being deleted needs no second request
		if st.ProvisioningStatus == "deleting" {
			continue
		}
		if sweep.ShouldSweep(resourceType, st.Name) {
			doomed = append(doomed, st)
		}
	}
	if err := iter.Err(); err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping object storage sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing object storages: %w", err)
	}

	for _, st := range doomed {
		log.Printf("[INFO] Deleting object storage: %s (%d)", st.Name, st.ID)
		if err := target.Client.Storage.ObjectStorages.DeleteAndPoll(ctx, st.ID); err != nil {
			sweep.RecordFailure(resourceType, st.Name, err)
			continue
		}
		sweep.RecordDeleted(resourceType, st.Name)
	}

	return nil
}
