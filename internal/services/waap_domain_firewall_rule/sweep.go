package waap_domain_firewall_rule

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/waap"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/G-Core/terraform-provider-gcore/internal/sweep"
)

func init() {
	resource.AddTestSweepers("gcore_waap_domain_firewall_rule", &resource.Sweeper{
		Name: "gcore_waap_domain_firewall_rule",
		F:    sweepWaapDomainFirewallRules,
	})
}

// sweepWaapDomainFirewallRules deletes test-prefixed firewall rules on the WAAP
// test domain (sweep.WaapDomainID), the only domain the acceptance tests use.
// Domains themselves are not managed by Terraform, so only the rules are swept.
func sweepWaapDomainFirewallRules(_ string) error {
	ctx := context.Background()
	// Resolve asserts GCORE_SWEEP_ACCOUNT_ID, so wrong credentials cannot sweep
	// another account's rules.
	target, err := sweep.Resolve(ctx)
	if err != nil {
		return err
	}
	client := target.Client

	domainID, err := sweep.WaapDomainID(ctx, client)
	if err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping WAAP firewall rule sweep: %s", err)
			return nil
		}
		return err
	}
	return sweepDomainFirewallRules(ctx, client, domainID)
}

// sweepDomainFirewallRules collects the test rules of one domain first and
// deletes them afterwards: deleting while offset-paging shifts the remaining
// rows down and skips some of them. A rule that is already gone counts as swept.
func sweepDomainFirewallRules(ctx context.Context, client *gcore.Client, domainID int64) error {
	var toDelete []waap.WaapFirewallRule
	rules := client.Waap.Domains.FirewallRules.ListAutoPaging(ctx, domainID, waap.DomainFirewallRuleListParams{})
	for rules.Next() {
		if rule := rules.Current(); sweep.ShouldSweep("gcore_waap_domain_firewall_rule", rule.Name) {
			toDelete = append(toDelete, rule)
		}
	}
	if err := rules.Err(); err != nil {
		return fmt.Errorf("error listing firewall rules for WAAP domain %d: %w", domainID, err)
	}

	var errs []error
	for _, rule := range toDelete {
		log.Printf("[INFO] Deleting WAAP firewall rule: %s (%d) on domain %d", rule.Name, rule.ID, domainID)
		err := client.Waap.Domains.FirewallRules.Delete(ctx, rule.ID, waap.DomainFirewallRuleDeleteParams{
			DomainID: domainID,
		})
		var apiErr *gcore.Error
		if err != nil && !(errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound) {
			errs = append(errs, fmt.Errorf("error deleting WAAP firewall rule %d on domain %d: %w", rule.ID, domainID, err))
		}
	}
	return errors.Join(errs...)
}
