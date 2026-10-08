package dns_network_mapping

import (
	"context"
	"fmt"
	"log"

	"github.com/G-Core/gcore-go/dns"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/G-Core/terraform-provider-gcore/internal/sweep"
)

func init() {
	resource.AddTestSweepers("gcore_dns_network_mapping", &resource.Sweeper{
		Name: "gcore_dns_network_mapping",
		F:    sweepDNSNetworkMappings,
	})
}

func sweepDNSNetworkMappings(_ string) error {
	const resourceType = "gcore_dns_network_mapping"
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
	// Account-global: this enumerates every network mapping in the account, so
	// the token match is the only thing bounding what gets deleted.
	var doomed []dns.DNSNetworkMapping
	iter := target.Client.DNS.NetworkMappings.ListAutoPaging(ctx, dns.NetworkMappingListParams{})
	for iter.Next() {
		mapping := iter.Current()
		sweep.RecordListed(resourceType, 1)
		if sweep.ShouldSweep(resourceType, mapping.Name) {
			doomed = append(doomed, mapping)
		}
	}
	if err := iter.Err(); err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping dns_network_mapping sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing dns network mappings: %w", err)
	}

	for _, mapping := range doomed {
		log.Printf("[INFO] Deleting dns_network_mapping: %s (ID: %d)", mapping.Name, mapping.ID)
		if _, err := target.Client.DNS.NetworkMappings.Delete(ctx, mapping.ID); err != nil {
			sweep.RecordFailure(resourceType, mapping.Name, err)
			continue
		}
		sweep.RecordDeleted(resourceType, mapping.Name)
	}

	return nil
}
