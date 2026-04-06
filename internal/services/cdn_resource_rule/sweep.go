package cdn_resource_rule

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/cdn"
	"github.com/G-Core/gcore-go/option"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/G-Core/terraform-provider-gcore/internal/sweep"
)

func init() {
	resource.AddTestSweepers("gcore_cdn_resource_rule", &resource.Sweeper{
		Name: "gcore_cdn_resource_rule",
		F:    sweepCDNResourceRules,
		// Sweep rules before CDN resources and origin groups
		Dependencies: []string{},
	})

	resource.AddTestSweepers("gcore_cdn_resource", &resource.Sweeper{
		Name: "gcore_cdn_resource",
		F:    sweepCDNResources,
		// Sweep CDN resources after their rules, but before origin groups
		Dependencies: []string{"gcore_cdn_resource_rule"},
	})
}

// sweepCDNResourceRules lists all CDN resources with test-prefixed cnames,
// then deletes any rules with test-prefixed names on those resources.
func sweepCDNResourceRules(_ string) error {
	if err := sweep.ValidateSweeperEnvironment(); err != nil {
		return err
	}

	apiKey := os.Getenv("GCORE_API_KEY")
	client := gcore.NewClient(option.WithAPIKey(apiKey))
	ctx := context.Background()

	cdnResources, err := client.CDN.CDNResources.List(ctx, cdn.CDNResourceListParams{})
	if err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping CDN resource rule sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing CDN resources: %w", err)
	}

	if cdnResources == nil {
		return nil
	}

	for _, cdnRes := range cdnResources.OfPlainList {
		if !isTestCDNResource(cdnRes) {
			continue
		}

		rules, err := client.CDN.CDNResources.Rules.List(ctx, cdnRes.ID, cdn.CDNResourceRuleListParams{})
		if err != nil {
			log.Printf("[ERROR] Failed to list rules for CDN resource %d: %s", cdnRes.ID, err)
			continue
		}
		if rules == nil {
			continue
		}

		for _, rule := range rules.Results {
			if !sweep.ShouldSweep("gcore_cdn_resource_rule", rule.Name) {
				continue
			}

			log.Printf("[INFO] Deleting CDN resource rule: %s (%d) on resource %d", rule.Name, rule.ID, cdnRes.ID)
			err := client.CDN.CDNResources.Rules.Delete(ctx, rule.ID, cdn.CDNResourceRuleDeleteParams{
				ResourceID: cdnRes.ID,
			})
			if err != nil {
				log.Printf("[ERROR] Failed to delete CDN resource rule %d: %s", rule.ID, err)
			}
		}
	}

	return nil
}

// sweepCDNResources lists all CDN resources and deletes those that were
// created by acceptance tests (identified by test-prefixed cnames).
func sweepCDNResources(_ string) error {
	if err := sweep.ValidateSweeperEnvironment(); err != nil {
		return err
	}

	apiKey := os.Getenv("GCORE_API_KEY")
	client := gcore.NewClient(option.WithAPIKey(apiKey))
	ctx := context.Background()

	cdnResources, err := client.CDN.CDNResources.List(ctx, cdn.CDNResourceListParams{})
	if err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping CDN resource sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing CDN resources: %w", err)
	}

	if cdnResources == nil {
		return nil
	}

	for _, cdnRes := range cdnResources.OfPlainList {
		if !isTestCDNResource(cdnRes) {
			continue
		}

		log.Printf("[INFO] Deleting CDN resource: %s (cname: %s, id: %d)", cdnRes.Name, cdnRes.Cname, cdnRes.ID)
		err := client.CDN.CDNResources.DeactivateAndDelete(ctx, cdnRes.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete CDN resource %d: %s", cdnRes.ID, err)
		}
	}

	return nil
}

// isTestCDNResource checks whether a CDN resource was created by acceptance
// tests. The fixture helper creates CDN resources with cnames like
// "tf-test-<random>.example.com", so we check the cname prefix.
func isTestCDNResource(res cdn.CDNResource) bool {
	return strings.HasPrefix(res.Cname, sweep.ResourcePrefix)
}
