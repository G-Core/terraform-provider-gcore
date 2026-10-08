package dns_zone_rrset

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/G-Core/gcore-go/dns"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/G-Core/terraform-provider-gcore/internal/sweep"
)

func init() {
	resource.AddTestSweepers("gcore_dns_zone_rrset", &resource.Sweeper{
		Name: "gcore_dns_zone_rrset",
		F:    sweepDNSZoneRrsets,
	})
}

func sweepDNSZoneRrsets(_ string) error {
	const resourceType = "gcore_dns_zone_rrset"
	defer sweep.Finish(resourceType)

	ctx := context.Background()
	target, err := sweep.Resolve(ctx)
	if err != nil {
		return err
	}

	// Collect the parent zones first, then their rrsets, then delete. Deleting
	// while iterating an offset-paginated collection shifts it, so the next
	// offset skips exactly as many items as were removed.
	var zones []string
	zoneIter := target.Client.DNS.Zones.ListAutoPaging(ctx, dns.ZoneListParams{})
	for zoneIter.Next() {
		zone := zoneIter.Current()
		// The gate is on the parent zone: an rrset is a test resource only
		// because the zone containing it is one.
		if sweep.ShouldSweep("gcore_dns_zone", zone.Name) {
			zones = append(zones, zone.Name)
		}
	}
	if err := zoneIter.Err(); err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping DNS zone rrset sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing DNS zones: %w", err)
	}

	type doomedRrset struct {
		zoneName string
		name     string
		rtype    string
	}

	// As for CDN rules: a zone whose rrsets cannot be listed fails the sweep,
	// but only after the other zones have been handled.
	var listErr error
	var doomed []doomedRrset
	for _, zoneName := range zones {
		rrsetIter := target.Client.DNS.Zones.Rrsets.ListAutoPaging(ctx, zoneName, dns.ZoneRrsetListParams{})
		for rrsetIter.Next() {
			rrset := rrsetIter.Current()
			sweep.RecordListed(resourceType, 1)
			if rrset.Type == "SOA" || rrset.Type == "NS" {
				sweep.RecordSkipped(resourceType, zoneName+"/"+rrset.Name+"/"+string(rrset.Type), "zone apex record, deleted with the zone")
				continue
			}
			doomed = append(doomed, doomedRrset{zoneName: zoneName, name: rrset.Name, rtype: string(rrset.Type)})
		}
		if err := rrsetIter.Err(); err != nil {
			if sweep.SkipSweepError(err) {
				log.Printf("[WARN] Skipping rrsets of zone %s: %s", zoneName, err)
				continue
			}
			listErr = errors.Join(listErr, fmt.Errorf("error listing rrsets for zone %s: %w", zoneName, err))
		}
	}

	for _, rrset := range doomed {
		log.Printf("[INFO] Deleting DNS rrset: %s/%s/%s", rrset.zoneName, rrset.name, rrset.rtype)
		if _, err := target.Client.DNS.Zones.Rrsets.Delete(ctx, rrset.rtype, dns.ZoneRrsetDeleteParams{
			ZoneName:  rrset.zoneName,
			RrsetName: rrset.name,
		}); err != nil {
			sweep.RecordFailure(resourceType, rrset.zoneName+"/"+rrset.name+"/"+rrset.rtype, err)
			continue
		}
		sweep.RecordDeleted(resourceType, rrset.zoneName+"/"+rrset.name+"/"+rrset.rtype)
	}

	return listErr
}
