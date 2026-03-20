package cdn_trusted_ca_certificate

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/cdn"
	"github.com/G-Core/gcore-go/option"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/G-Core/terraform-provider-gcore/internal/sweep"
)

func init() {
	resource.AddTestSweepers("gcore_cdn_trusted_ca_certificate", &resource.Sweeper{
		Name: "gcore_cdn_trusted_ca_certificate",
		F:    sweepCDNTrustedCaCertificates,
	})
}

func sweepCDNTrustedCaCertificates(_ string) error {
	if err := sweep.ValidateSweeperEnvironment(); err != nil {
		return err
	}

	apiKey := os.Getenv("GCORE_API_KEY")
	client := gcore.NewClient(option.WithAPIKey(apiKey))

	ctx := context.Background()

	certs, err := client.CDN.TrustedCaCertificates.List(ctx, cdn.TrustedCaCertificateListParams{})
	if err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping CDN trusted CA certificate sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing CDN trusted CA certificates: %w", err)
	}

	for _, cert := range *certs {
		if !sweep.ShouldSweep("gcore_cdn_trusted_ca_certificate", cert.Name) {
			continue
		}

		log.Printf("[INFO] Deleting CDN trusted CA certificate: %s (%d)", cert.Name, cert.ID)
		err := client.CDN.TrustedCaCertificates.Delete(ctx, cert.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete CDN trusted CA certificate %s (%d): %s", cert.Name, cert.ID, err)
		}
	}

	return nil
}
