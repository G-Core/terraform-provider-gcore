package cloud_inference_secret

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
	resource.AddTestSweepers("gcore_cloud_inference_secret", &resource.Sweeper{
		Name: "gcore_cloud_inference_secret",
		F:    sweepCloudInferenceSecrets,
	})
}

func sweepCloudInferenceSecrets(_ string) error {
	const resourceType = "gcore_cloud_inference_secret"
	defer sweep.Finish(resourceType)

	ctx := context.Background()
	target, err := sweep.Resolve(ctx)
	if err != nil {
		return err
	}

	// Collect first, delete second. Deleting while iterating an
	// offset-paginated collection shifts it, so the next offset skips exactly
	// as many items as were removed.
	var doomed []cloud.InferenceSecret
	iter := target.Client.Cloud.Inference.Secrets.ListAutoPaging(ctx, cloud.InferenceSecretListParams{
		ProjectID: param.NewOpt(target.ProjectID),
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
			log.Printf("[WARN] Skipping inference secret sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing inference secrets: %w", err)
	}

	for _, secret := range doomed {
		log.Printf("[INFO] Deleting inference secret: %s", secret.Name)
		if err := target.Client.Cloud.Inference.Secrets.Delete(ctx, secret.Name, cloud.InferenceSecretDeleteParams{
			ProjectID: param.NewOpt(target.ProjectID),
		}); err != nil {
			sweep.RecordFailure(resourceType, secret.Name, err)
			continue
		}
		sweep.RecordDeleted(resourceType, secret.Name)
	}

	return nil
}
