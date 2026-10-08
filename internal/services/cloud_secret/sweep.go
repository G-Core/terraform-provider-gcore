package cloud_secret

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
	resource.AddTestSweepers("gcore_cloud_secret", &resource.Sweeper{
		Name: "gcore_cloud_secret",
		F:    sweepCloudSecrets,
	})
}

func sweepCloudSecrets(_ string) error {
	const resourceType = "gcore_cloud_secret"
	defer sweep.Finish(resourceType)

	ctx := context.Background()
	target, err := sweep.Resolve(ctx)
	if err != nil {
		return err
	}

	// Collect first, delete second. Deleting while iterating an
	// offset-paginated collection shifts it, so the next offset skips exactly
	// as many items as were removed.
	var doomed []cloud.Secret
	iter := target.Client.Cloud.Secrets.ListAutoPaging(ctx, cloud.SecretListParams{
		ProjectID: param.NewOpt(target.ProjectID),
		RegionID:  param.NewOpt(target.RegionID),
	})
	for iter.Next() {
		secret := iter.Current()
		sweep.RecordListed(resourceType, 1)
		if sweep.ShouldSweep(resourceType, secret.Name) {
			doomed = append(doomed, secret)
		}
	}
	if err := iter.Err(); err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping secret sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing secrets: %w", err)
	}

	for _, secret := range doomed {
		log.Printf("[INFO] Deleting secret: %s (%s)", secret.Name, secret.ID)
		if err := target.Client.Cloud.Secrets.DeleteAndPoll(ctx, secret.ID, cloud.SecretDeleteParams{
			ProjectID: param.NewOpt(target.ProjectID),
			RegionID:  param.NewOpt(target.RegionID),
		}); err != nil {
			sweep.RecordFailure(resourceType, secret.Name, err)
			continue
		}
		sweep.RecordDeleted(resourceType, secret.Name)
	}

	return nil
}
