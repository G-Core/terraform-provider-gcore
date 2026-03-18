package dns_network_mapping

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/dns"
	"github.com/G-Core/gcore-go/option"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("gcore_dns_network_mapping", &resource.Sweeper{
		Name: "gcore_dns_network_mapping",
		F:    sweepDNSNetworkMappings,
	})
}

func sweepDNSNetworkMappings(_ string) error {
	apiKey := os.Getenv("GCORE_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("GCORE_API_KEY must be set for sweepers")
	}

	client := gcore.NewClient(option.WithAPIKey(apiKey))
	ctx := context.Background()

	resp, err := client.DNS.NetworkMappings.List(ctx, dns.NetworkMappingListParams{})
	if err != nil {
		log.Printf("[WARN] Skipping dns_network_mapping sweep: %s", err)
		return nil
	}

	for _, mapping := range resp.NetworkMappings {
		if !strings.HasPrefix(mapping.Name, "tf-test") && !strings.HasPrefix(mapping.Name, "tf_test") {
			log.Printf("[DEBUG] Skipping dns_network_mapping: %s (not a test resource)", mapping.Name)
			continue
		}

		log.Printf("[INFO] Deleting dns_network_mapping: %s (ID: %d)", mapping.Name, mapping.ID)
		_, err := client.DNS.NetworkMappings.Delete(ctx, mapping.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete dns_network_mapping %s: %s", mapping.Name, err)
		}
	}

	return nil
}
