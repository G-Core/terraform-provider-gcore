package dns_zone_rrset

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/dns"
	"github.com/G-Core/gcore-go/option"
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
	if err := sweep.ValidateSweeperEnvironment(); err != nil {
		return err
	}

	apiKey := os.Getenv("GCORE_API_KEY")

	client := gcore.NewClient(option.WithAPIKey(apiKey))
	ctx := context.Background()

	zones, err := client.DNS.Zones.List(ctx, dns.ZoneListParams{})
	if err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping DNS zone rrset sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing DNS zones: %w", err)
	}

	for _, zone := range zones.Zones {
		if !sweep.ShouldSweep("gcore_dns_zone", zone.Name) {
			continue
		}

		rrsets, err := client.DNS.Zones.Rrsets.List(ctx, zone.Name, dns.ZoneRrsetListParams{})
		if err != nil {
			log.Printf("[WARN] Failed to list rrsets for zone %s: %s", zone.Name, err)
			continue
		}

		for _, rrset := range rrsets.Rrsets {
			if rrset.Type == "SOA" || rrset.Type == "NS" {
				continue
			}

			log.Printf("[INFO] Deleting DNS rrset: %s/%s/%s", zone.Name, rrset.Name, rrset.Type)
			_, err := client.DNS.Zones.Rrsets.Delete(ctx, string(rrset.Type), dns.ZoneRrsetDeleteParams{
				ZoneName:  zone.Name,
				RrsetName: rrset.Name,
			})
			if err != nil {
				log.Printf("[ERROR] Failed to delete rrset %s/%s/%s: %s", zone.Name, rrset.Name, rrset.Type, err)
			}
		}
	}

	return nil
}
