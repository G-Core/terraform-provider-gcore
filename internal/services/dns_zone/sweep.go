package dns_zone

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/dns"
	"github.com/G-Core/gcore-go/option"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/stainless-sdks/gcore-terraform/internal/sweep"
)

func init() {
	resource.AddTestSweepers("gcore_dns_zone", &resource.Sweeper{
		Name: "gcore_dns_zone",
		F:    sweepDNSZones,
	})
}

func sweepDNSZones(_ string) error {
	if err := sweep.ValidateSweeperEnvironment(); err != nil {
		return err
	}

	apiKey := os.Getenv("GCORE_API_KEY")

	client := gcore.NewClient(
		option.WithAPIKey(apiKey),
	)

	ctx := context.Background()

	page, err := client.DNS.Zones.List(ctx, dns.ZoneListParams{})
	if err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping DNS zone sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing DNS zones: %w", err)
	}

	for _, zone := range page.Zones {
		if !sweep.ShouldSweep("gcore_dns_zone", zone.Name) {
			continue
		}

		log.Printf("[INFO] Deleting DNS zone: %s", zone.Name)
		_, err := client.DNS.Zones.Delete(ctx, zone.Name)
		if err != nil {
			log.Printf("[ERROR] Failed to delete DNS zone %s: %s", zone.Name, err)
		}
	}

	return nil
}
