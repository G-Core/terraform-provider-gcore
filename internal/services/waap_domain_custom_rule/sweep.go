package waap_domain_custom_rule

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/G-Core/gcore-go/waap"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/G-Core/terraform-provider-gcore/internal/sweep"
)

func init() {
	resource.AddTestSweepers("gcore_waap_domain_custom_rule", &resource.Sweeper{
		Name: "gcore_waap_domain_custom_rule",
		F:    sweepWaapDomainCustomRules,
	})
}

// sweepWaapDomainCustomRules deletes the custom rules of the run being swept on
// the WAAP test domain (sweep.WaapDomainID), the only domain the acceptance
// tests use. The domain itself is not a test resource and is left alone. A plan
// can allow as few as one custom rule per domain, so a leaked rule blocks every
// later test run until it is swept.
func sweepWaapDomainCustomRules(_ string) error {
	const resourceType = "gcore_waap_domain_custom_rule"
	defer sweep.Finish(resourceType)

	ctx := context.Background()
	target, err := sweep.Resolve(ctx)
	if err != nil {
		return err
	}
	client := target.Client

	domainID, err := sweep.WaapDomainID(ctx, client)
	if err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping WAAP custom rule sweep: %s", err)
			return nil
		}
		return err
	}

	type doomedRule struct {
		id   int64
		name string
	}

	// Collect first, delete second: deleting while offset-paging shifts later
	// rules onto pages already read, so they would be skipped. A listing that
	// fails part-way is returned once the deletes are done.
	var errs []error
	var doomed []doomedRule
	rules := client.Waap.Domains.CustomRules.ListAutoPaging(ctx, domainID, waap.DomainCustomRuleListParams{})
	for rules.Next() {
		rule := rules.Current()
		sweep.RecordListed(resourceType, 1)
		if sweep.ShouldSweep(resourceType, rule.Name) {
			doomed = append(doomed, doomedRule{id: rule.ID, name: rule.Name})
		}
	}
	if err := rules.Err(); err != nil && !sweep.SkipSweepError(err) {
		errs = append(errs, fmt.Errorf("error listing custom rules for WAAP domain %d: %w", domainID, err))
	}

	for _, rule := range doomed {
		log.Printf("[INFO] Deleting WAAP custom rule: %s (%d) on domain %d", rule.name, rule.id, domainID)
		err := client.Waap.Domains.CustomRules.Delete(ctx, rule.id, waap.DomainCustomRuleDeleteParams{
			DomainID: domainID,
		})
		if err != nil {
			sweep.RecordFailure(resourceType, rule.name, err)
			errs = append(errs, fmt.Errorf("error deleting WAAP custom rule %d on domain %d: %w", rule.id, domainID, err))
			continue
		}
		sweep.RecordDeleted(resourceType, rule.name)
	}

	return errors.Join(errs...)
}
