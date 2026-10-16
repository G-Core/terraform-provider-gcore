package waap_domain_advanced_rule

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
	resource.AddTestSweepers("gcore_waap_domain_advanced_rule", &resource.Sweeper{
		Name: "gcore_waap_domain_advanced_rule",
		F:    sweepWaapDomainAdvancedRules,
	})
}

// sweepWaapDomainAdvancedRules deletes the advanced rules of the run being
// swept on the WAAP test domain (sweep.WaapDomainID), the only domain the
// acceptance tests use. The domain itself is not a test resource and is left
// alone.
func sweepWaapDomainAdvancedRules(_ string) error {
	const resourceType = "gcore_waap_domain_advanced_rule"
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
			log.Printf("[WARN] Skipping WAAP advanced rule sweep: %s", err)
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
	rules := client.Waap.Domains.AdvancedRules.ListAutoPaging(ctx, domainID, waap.DomainAdvancedRuleListParams{})
	for rules.Next() {
		rule := rules.Current()
		sweep.RecordListed(resourceType, 1)
		if sweep.ShouldSweep(resourceType, rule.Name) {
			doomed = append(doomed, doomedRule{id: rule.ID, name: rule.Name})
		}
	}
	if err := rules.Err(); err != nil && !sweep.SkipSweepError(err) {
		errs = append(errs, fmt.Errorf("error listing advanced rules for WAAP domain %d: %w", domainID, err))
	}

	for _, rule := range doomed {
		log.Printf("[INFO] Deleting WAAP advanced rule: %s (%d) on domain %d", rule.name, rule.id, domainID)
		err := client.Waap.Domains.AdvancedRules.Delete(ctx, rule.id, waap.DomainAdvancedRuleDeleteParams{
			DomainID: domainID,
		})
		if err != nil {
			sweep.RecordFailure(resourceType, rule.name, err)
			errs = append(errs, fmt.Errorf("error deleting WAAP advanced rule %d on domain %d: %w", rule.id, domainID, err))
			continue
		}
		sweep.RecordDeleted(resourceType, rule.name)
	}

	return errors.Join(errs...)
}
