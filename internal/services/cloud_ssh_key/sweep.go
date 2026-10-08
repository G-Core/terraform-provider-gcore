package cloud_ssh_key

import (
	"context"
	"fmt"
	"log"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/G-Core/terraform-provider-gcore/internal/sweep"
)

func init() {
	resource.AddTestSweepers("gcore_cloud_ssh_key", &resource.Sweeper{
		Name: "gcore_cloud_ssh_key",
		F:    sweepCloudSSHKeys,
	})
}

func sweepCloudSSHKeys(_ string) error {
	const resourceType = "gcore_cloud_ssh_key"
	defer sweep.Finish(resourceType)

	ctx := context.Background()
	target, err := sweep.Resolve(ctx)
	if err != nil {
		return err
	}

	// Collect first, delete second. Deleting while iterating an
	// offset-paginated collection shifts it, so the next offset skips exactly
	// as many items as were removed.
	var doomed []cloud.SSHKey
	iter := target.Client.Cloud.SSHKeys.ListAutoPaging(ctx, cloud.SSHKeyListParams{
		ProjectID: param.NewOpt(target.ProjectID),
	})
	for iter.Next() {
		key := iter.Current()
		sweep.RecordListed(resourceType, 1)
		if sweep.ShouldSweep(resourceType, key.Name) {
			doomed = append(doomed, key)
		}
	}
	if err := iter.Err(); err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping SSH key sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing SSH keys: %w", err)
	}

	for _, key := range doomed {
		log.Printf("[INFO] Deleting SSH key: %s (%s)", key.Name, key.ID)
		err := target.Client.Cloud.SSHKeys.Delete(ctx, key.ID, cloud.SSHKeyDeleteParams{
			ProjectID: param.NewOpt(target.ProjectID),
		})
		if err != nil {
			sweep.RecordFailure(resourceType, key.Name, err)
			continue
		}
		sweep.RecordDeleted(resourceType, key.Name)
	}

	return nil
}
