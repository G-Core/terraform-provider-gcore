package cdn_origin_group

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
	resource.AddTestSweepers("gcore_cdn_origin_group", &resource.Sweeper{
		Name: "gcore_cdn_origin_group",
		F:    sweepCDNOriginGroups,
	})
}

func sweepCDNOriginGroups(_ string) error {
	if err := sweep.ValidateSweeperEnvironment(); err != nil {
		return err
	}

	apiKey := os.Getenv("GCORE_API_KEY")
	client := gcore.NewClient(option.WithAPIKey(apiKey))

	ctx := context.Background()

	groups, err := client.CDN.OriginGroups.List(ctx, cdn.OriginGroupListParams{})
	if err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping CDN origin group sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing CDN origin groups: %w", err)
	}

	if groups == nil {
		return nil
	}

	for _, group := range groups.OfPlainList {
		if !sweep.ShouldSweep("gcore_cdn_origin_group", group.Name) {
			continue
		}

		// Skip groups with related resources - cannot delete if CDN resources are using it
		if group.HasRelatedResources {
			log.Printf("[INFO] Skipping CDN origin group with related resources: %s (%d)", group.Name, group.ID)
			continue
		}

		log.Printf("[INFO] Deleting CDN origin group: %s (%d)", group.Name, group.ID)
		err := client.CDN.OriginGroups.Delete(ctx, group.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete CDN origin group %s (%d): %s", group.Name, group.ID, err)
		}
	}

	return nil
}
