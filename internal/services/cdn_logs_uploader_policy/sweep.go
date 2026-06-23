package cdn_logs_uploader_policy

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
	resource.AddTestSweepers("gcore_cdn_logs_uploader_policy", &resource.Sweeper{
		Name: "gcore_cdn_logs_uploader_policy",
		F:    sweepCDNLogsUploaderPolicies,
	})
}

func sweepCDNLogsUploaderPolicies(_ string) error {
	if err := sweep.ValidateSweeperEnvironment(); err != nil {
		return err
	}

	apiKey := os.Getenv("GCORE_API_KEY")
	client := gcore.NewClient(option.WithAPIKey(apiKey))

	ctx := context.Background()

	policies, err := client.CDN.LogsUploader.Policies.List(ctx, cdn.LogsUploaderPolicyListParams{})
	if err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping CDN logs uploader policy sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing CDN logs uploader policies: %w", err)
	}

	allPolicies := policies.OfPlainList
	if len(allPolicies) == 0 {
		allPolicies = policies.Results
	}
	for _, policy := range allPolicies {
		if !sweep.ShouldSweep("gcore_cdn_logs_uploader_policy", policy.Name) {
			continue
		}

		log.Printf("[INFO] Deleting CDN logs uploader policy: %s (%d)", policy.Name, policy.ID)
		err := client.CDN.LogsUploader.Policies.Delete(ctx, policy.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete CDN logs uploader policy %s (%d): %s", policy.Name, policy.ID, err)
		}
	}

	return nil
}
