package cdn_trusted_ca_certificate

import (
	"context"
	"fmt"
	"log"

	"github.com/G-Core/gcore-go/cdn"
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
	const resourceType = "gcore_cdn_trusted_ca_certificate"
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
	// Account-global: this enumerates every trusted CA certificate in the
	// account, so the token match is the only thing bounding what gets deleted.
	var doomed []cdn.CaCertificate
	iter := target.Client.CDN.TrustedCaCertificates.ListAutoPaging(ctx, cdn.TrustedCaCertificateListParams{})
	for iter.Next() {
		cert := iter.Current()
		sweep.RecordListed(resourceType, 1)
		if sweep.ShouldSweep(resourceType, cert.Name) {
			doomed = append(doomed, cert)
		}
	}
	if err := iter.Err(); err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping CDN trusted CA certificate sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing CDN trusted CA certificates: %w", err)
	}

	for _, cert := range doomed {
		log.Printf("[INFO] Deleting CDN trusted CA certificate: %s (%d)", cert.Name, cert.ID)
		if err := target.Client.CDN.TrustedCaCertificates.Delete(ctx, cert.ID); err != nil {
			sweep.RecordFailure(resourceType, cert.Name, err)
			continue
		}
		sweep.RecordDeleted(resourceType, cert.Name)
	}

	return nil
}
