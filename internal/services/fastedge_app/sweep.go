package fastedge_app

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/fastedge"
	"github.com/G-Core/gcore-go/option"
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
	if err := sweep.ValidateSweeperEnvironment(); err != nil {
		return err
	}

	apiKey := os.Getenv("GCORE_API_KEY")
	client := gcore.NewClient(option.WithAPIKey(apiKey))

	ctx := context.Background()

	params := fastedge.AppListParams{}
	page, err := client.Fastedge.Apps.List(ctx, params)
	if err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping fastedge app sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing fastedge apps: %w", err)
	}

	for page != nil && len(page.Apps) > 0 {
		for _, app := range page.Apps {
			if !sweep.ShouldSweep("gcore_fastedge_app", app.Name) {
				continue
			}
			log.Printf("[INFO] Deleting fastedge app: %s (%d)", app.Name, app.ID)
			err := client.Fastedge.Apps.Delete(ctx, app.ID)
			if err != nil {
				log.Printf("[ERROR] Failed to delete fastedge app %s: %s", app.Name, err)
			}
		}

		page, err = page.GetNextPage()
		if err != nil {
			log.Printf("[ERROR] Failed to get next page of apps: %s", err)
			break
		}
	}

	return nil
}
