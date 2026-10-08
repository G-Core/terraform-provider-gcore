package cdn_resource_rule

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/G-Core/gcore-go/cdn"
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
// then deletes the rules on them that belong to the run being swept. The cname
// check only narrows which resources to look inside; ShouldSweep on each rule
// name decides what is deleted.
func sweepCDNResourceRules(_ string) error {
	const resourceType = "gcore_cdn_resource_rule"
	defer sweep.Finish(resourceType)

	ctx := context.Background()
	target, err := sweep.Resolve(ctx)
	if err != nil {
		return err
	}

	// Collect the parent resources first, then their rules, then delete.
	// Deleting while iterating an offset-paginated collection shifts it, so the
	// next offset skips exactly as many items as were removed.
	var parents []cdn.CDNResource
	resIter := target.Client.CDN.CDNResources.ListAutoPaging(ctx, cdn.CDNResourceListParams{})
	for resIter.Next() {
		cdnRes := resIter.Current()
		if isTestCDNResource(cdnRes) {
			parents = append(parents, cdnRes)
		}
	}
	if err := resIter.Err(); err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping CDN resource rule sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing CDN resources: %w", err)
	}

	type doomedRule struct {
		resourceID int64
		id         int64
		name       string
	}

	// A rule listing that fails part-way must not read as a clean sweep. Keep
	// going so one resource's rules do not abandon the rest, but return the
	// error once the deletes are done.
	var listErr error
	var doomed []doomedRule
	for _, cdnRes := range parents {
		ruleIter := target.Client.CDN.CDNResources.Rules.ListAutoPaging(ctx, cdnRes.ID, cdn.CDNResourceRuleListParams{})
		for ruleIter.Next() {
			rule := ruleIter.Current()
			sweep.RecordListed(resourceType, 1)
			if sweep.ShouldSweep(resourceType, rule.Name) {
				doomed = append(doomed, doomedRule{resourceID: cdnRes.ID, id: rule.ID, name: rule.Name})
			}
		}
		if err := ruleIter.Err(); err != nil {
			if sweep.SkipSweepError(err) {
				log.Printf("[WARN] Skipping rules of CDN resource %d: %s", cdnRes.ID, err)
				continue
			}
			listErr = errors.Join(listErr, fmt.Errorf("error listing rules for CDN resource %d: %w", cdnRes.ID, err))
		}
	}

	for _, rule := range doomed {
		log.Printf("[INFO] Deleting CDN resource rule: %s (%d) on resource %d", rule.name, rule.id, rule.resourceID)
		if err := target.Client.CDN.CDNResources.Rules.Delete(ctx, rule.id, cdn.CDNResourceRuleDeleteParams{
			ResourceID: rule.resourceID,
		}); err != nil {
			sweep.RecordFailure(resourceType, rule.name, err)
			continue
		}
		sweep.RecordDeleted(resourceType, rule.name)
	}

	return listErr
}

// sweepCDNResources lists all CDN resources and deletes those created by the
// run being swept. Test cnames are "<RandomName>.example.com", so they carry
// the run token and go through ShouldSweep like any other name.
func sweepCDNResources(_ string) error {
	const resourceType = "gcore_cdn_resource"
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
	// Account-global: this enumerates every CDN resource in the account, so the
	// cname's run token is the only thing bounding what gets deleted.
	var doomed []cdn.CDNResource
	iter := target.Client.CDN.CDNResources.ListAutoPaging(ctx, cdn.CDNResourceListParams{})
	for iter.Next() {
		cdnRes := iter.Current()
		sweep.RecordListed(resourceType, 1)
		if sweep.ShouldSweep(resourceType, cdnRes.Cname) {
			doomed = append(doomed, cdnRes)
		}
	}
	if err := iter.Err(); err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping CDN resource sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing CDN resources: %w", err)
	}

	for _, cdnRes := range doomed {
		log.Printf("[INFO] Deleting CDN resource: %s (cname: %s, id: %d)", cdnRes.Name, cdnRes.Cname, cdnRes.ID)
		if err := target.Client.CDN.CDNResources.DeactivateAndDelete(ctx, cdnRes.ID); err != nil {
			sweep.RecordFailure(resourceType, cdnRes.Cname, err)
			continue
		}
		sweep.RecordDeleted(resourceType, cdnRes.Cname)
	}

	return nil
}

// isTestCDNResource checks whether a CDN resource was created by acceptance
// tests, whatever run it belongs to. It must not gate deletion of the resource
// itself; use sweep.ShouldSweep for that.
func isTestCDNResource(res cdn.CDNResource) bool {
	return sweep.IsTestResource(res.Cname)
}
