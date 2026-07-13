package cdn_logs_uploader_target

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/cdn"
	"github.com/G-Core/gcore-go/option"
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
	if err := sweep.ValidateSweeperEnvironment(); err != nil {
		return err
	}

	apiKey := os.Getenv("GCORE_API_KEY")
	client := gcore.NewClient(option.WithAPIKey(apiKey))

	ctx := context.Background()

	targets, err := client.CDN.LogsUploader.Targets.List(ctx, cdn.LogsUploaderTargetListParams{})
	if err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping CDN logs uploader target sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing CDN logs uploader targets: %w", err)
	}

	// CDN list endpoints are now offset pagers; the result exposes .Results
	// (the {count,results} envelope), not the plain-array union.
	for _, target := range targets.Results {
		if !sweep.ShouldSweep("gcore_cdn_logs_uploader_target", target.Name) {
			continue
		}

		log.Printf("[INFO] Deleting CDN logs uploader target: %s (%d)", target.Name, target.ID)
		err := client.CDN.LogsUploader.Targets.Delete(ctx, target.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete CDN logs uploader target %s (%d): %s", target.Name, target.ID, err)
		}
	}

	return nil
}
