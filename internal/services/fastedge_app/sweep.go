package fastedge_app

import (
	"context"
	"fmt"
	"log"

	"github.com/G-Core/gcore-go/fastedge"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/G-Core/terraform-provider-gcore/internal/sweep"
)

func init() {
	resource.AddTestSweepers("gcore_fastedge_app", &resource.Sweeper{
		Name: "gcore_fastedge_app",
		F:    sweepFastedgeApps,
		// Note: apps should be deleted before binaries. When a fastedge_binary
		// sweeper is created, it should declare this sweeper as a dependency.
	})
}

func sweepFastedgeApps(_ string) error {
	const resourceType = "gcore_fastedge_app"
	defer sweep.Finish(resourceType)

	ctx := context.Background()
	target, err := sweep.Resolve(ctx)
	if err != nil {
		return err
	}

	// Collect first, delete second. This previously deleted while advancing the
	// pager, and each delete shifts an offset-paginated collection, so the next
	// offset skipped exactly as many items as had been removed.
	//
	// Account-global: this enumerates every app in the account, so the token
	// match is the only thing bounding what gets deleted.
	var doomed []fastedge.AppShort
	iter := target.Client.Fastedge.Apps.ListAutoPaging(ctx, fastedge.AppListParams{})
	for iter.Next() {
		app := iter.Current()
		sweep.RecordListed(resourceType, 1)
		if sweep.ShouldSweep(resourceType, app.Name) {
			doomed = append(doomed, app)
		}
	}
	if err := iter.Err(); err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping fastedge app sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing fastedge apps: %w", err)
	}

	for _, app := range doomed {
		log.Printf("[INFO] Deleting fastedge app: %s (%d)", app.Name, app.ID)
		if err := target.Client.Fastedge.Apps.Delete(ctx, app.ID); err != nil {
			sweep.RecordFailure(resourceType, app.Name, err)
			continue
		}
		sweep.RecordDeleted(resourceType, app.Name)
	}

	return nil
}
