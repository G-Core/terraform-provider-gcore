package cdn_logs_uploader_policy

import (
	"context"
	"fmt"
	"log"

	"github.com/G-Core/gcore-go/cdn"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/G-Core/terraform-provider-gcore/internal/sweep"
)

func init() {
	resource.AddTestSweepers("gcore_cdn_logs_uploader_policy", &resource.Sweeper{
		Name: "gcore_cdn_logs_uploader_policy",
		F:    sweepCDNLogsUploaderPolicies,
	})
}

func sweepCDNLogsUploaderPolicies(_ string) error {
	const resourceType = "gcore_cdn_logs_uploader_policy"
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
	// Account-global: this enumerates every policy in the account, so the token
	// match is the only thing bounding what gets deleted.
	var doomed []cdn.LogsUploaderPolicy
	iter := target.Client.CDN.LogsUploader.Policies.ListAutoPaging(ctx, cdn.LogsUploaderPolicyListParams{})
	for iter.Next() {
		policy := iter.Current()
		sweep.RecordListed(resourceType, 1)
		if sweep.ShouldSweep(resourceType, policy.Name) {
			doomed = append(doomed, policy)
		}
	}
	if err := iter.Err(); err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping CDN logs uploader policy sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing CDN logs uploader policies: %w", err)
	}

	for _, policy := range doomed {
		log.Printf("[INFO] Deleting CDN logs uploader policy: %s (%d)", policy.Name, policy.ID)
		if err := target.Client.CDN.LogsUploader.Policies.Delete(ctx, policy.ID); err != nil {
			sweep.RecordFailure(resourceType, policy.Name, err)
			continue
		}
		sweep.RecordDeleted(resourceType, policy.Name)
	}

	return nil
}
