package cdn_logs_uploader_target

import (
	"context"
	"fmt"
	"log"

	"github.com/G-Core/gcore-go/cdn"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/G-Core/terraform-provider-gcore/internal/sweep"
)

func init() {
	resource.AddTestSweepers("gcore_cdn_logs_uploader_target", &resource.Sweeper{
		Name: "gcore_cdn_logs_uploader_target",
		F:    sweepCDNLogsUploaderTargets,
	})
}

func sweepCDNLogsUploaderTargets(_ string) error {
	const resourceType = "gcore_cdn_logs_uploader_target"
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
	// Account-global: this enumerates every target in the account, so the token
	// match is the only thing bounding what gets deleted.
	var doomed []cdn.LogsUploaderTargetUnion
	iter := target.Client.CDN.LogsUploader.Targets.ListAutoPaging(ctx, cdn.LogsUploaderTargetListParams{})
	for iter.Next() {
		uploadTarget := iter.Current()
		sweep.RecordListed(resourceType, 1)
		if sweep.ShouldSweep(resourceType, uploadTarget.Name) {
			doomed = append(doomed, uploadTarget)
		}
	}
	if err := iter.Err(); err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping CDN logs uploader target sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing CDN logs uploader targets: %w", err)
	}

	for _, uploadTarget := range doomed {
		log.Printf("[INFO] Deleting CDN logs uploader target: %s (%d)", uploadTarget.Name, uploadTarget.ID)
		if err := target.Client.CDN.LogsUploader.Targets.Delete(ctx, uploadTarget.ID); err != nil {
			sweep.RecordFailure(resourceType, uploadTarget.Name, err)
			continue
		}
		sweep.RecordDeleted(resourceType, uploadTarget.Name)
	}

	return nil
}
