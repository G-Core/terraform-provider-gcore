package dns_zone

import (
	"context"
	"fmt"
	"log"

	"github.com/G-Core/gcore-go/dns"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/G-Core/terraform-provider-gcore/internal/sweep"
)

func init() {
	resource.AddTestSweepers("gcore_dns_zone", &resource.Sweeper{
		Name:         "gcore_dns_zone",
		F:            sweepDNSZones,
		Dependencies: []string{"gcore_dns_zone_rrset"},
	})
}

func sweepDNSZones(_ string) error {
	const resourceType = "gcore_dns_zone"
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
	// DNS zones are account-global: this enumerates every zone in the account,
	// not just a project's, so the token match is the only thing bounding what
	// gets deleted.
	var doomed []dns.ZoneListResponse
	iter := target.Client.DNS.Zones.ListAutoPaging(ctx, dns.ZoneListParams{})
	for iter.Next() {
		zone := iter.Current()
		sweep.RecordListed(resourceType, 1)
		if sweep.ShouldSweep(resourceType, zone.Name) {
			doomed = append(doomed, zone)
		}
	}
	if err := iter.Err(); err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping DNS zone sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing DNS zones: %w", err)
	}

	for _, zone := range doomed {
		log.Printf("[INFO] Deleting DNS zone: %s", zone.Name)
		if _, err := target.Client.DNS.Zones.Delete(ctx, zone.Name); err != nil {
			sweep.RecordFailure(resourceType, zone.Name, err)
			continue
		}
		sweep.RecordDeleted(resourceType, zone.Name)
	}

	return nil
}
