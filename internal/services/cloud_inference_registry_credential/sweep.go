package cloud_inference_registry_credential

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
	resource.AddTestSweepers("gcore_cloud_inference_registry_credential", &resource.Sweeper{
		Name: "gcore_cloud_inference_registry_credential",
		F:    sweepCloudInferenceRegistryCredentials,
	})
}

func sweepCloudInferenceRegistryCredentials(_ string) error {
	const resourceType = "gcore_cloud_inference_registry_credential"
	defer sweep.Finish(resourceType)

	ctx := context.Background()
	target, err := sweep.Resolve(ctx)
	if err != nil {
		return err
	}

	// Collect first, delete second. Deleting while iterating an
	// offset-paginated collection shifts it, so the next offset skips exactly
	// as many items as were removed.
	var doomed []cloud.InferenceRegistryCredentials
	iter := target.Client.Cloud.Inference.RegistryCredentials.ListAutoPaging(ctx, cloud.InferenceRegistryCredentialListParams{
		ProjectID: param.NewOpt(target.ProjectID),
	})
	for iter.Next() {
		cred := iter.Current()
		sweep.RecordListed(resourceType, 1)
		if sweep.ShouldSweep(resourceType, cred.Name) {
			doomed = append(doomed, cred)
		}
	}
	if err := iter.Err(); err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping inference registry credential sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing inference registry credentials: %w", err)
	}

	for _, cred := range doomed {
		log.Printf("[INFO] Deleting inference registry credential: %s", cred.Name)
		if err := target.Client.Cloud.Inference.RegistryCredentials.Delete(ctx, cred.Name, cloud.InferenceRegistryCredentialDeleteParams{
			ProjectID: param.NewOpt(target.ProjectID),
		}); err != nil {
			sweep.RecordFailure(resourceType, cred.Name, err)
			continue
		}
		sweep.RecordDeleted(resourceType, cred.Name)
	}

	return nil
}
