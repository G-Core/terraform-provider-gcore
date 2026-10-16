package sweep

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/gcore-go/waap"
)

// WaapTestDomainName is the WAAP domain the WAAP acceptance tests attach rules
// to and the WAAP sweepers clean. A WAAP domain needs a CDN resource with WAAP
// enabled, so it lives in the test account instead of being created per run.
const WaapTestDomainName = "terraform-acctest-waap.example.com"

// EnvWaapDomainID overrides the lookup by name, for an account whose test
// domain has another name.
const EnvWaapDomainID = "GCORE_WAAP_DOMAIN_ID"

// WaapDomainID returns the WAAP test domain: EnvWaapDomainID when set,
// otherwise the domain named WaapTestDomainName. It never falls back to another
// domain of the account, so tests and sweepers cannot touch a domain that is
// not meant for them.
func WaapDomainID(ctx context.Context, client *gcore.Client) (int64, error) {
	if v := os.Getenv(EnvWaapDomainID); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid %s %q: %w", EnvWaapDomainID, v, err)
		}
		return id, nil
	}

	// The name filter may match partially, so compare the name exactly.
	domains := client.Waap.Domains.ListAutoPaging(ctx, waap.DomainListParams{
		Name: param.NewOpt(WaapTestDomainName),
	})
	for domains.Next() {
		if d := domains.Current(); d.Name == WaapTestDomainName {
			return d.ID, nil
		}
	}
	if err := domains.Err(); err != nil {
		return 0, fmt.Errorf("error listing WAAP domains: %w", err)
	}
	return 0, fmt.Errorf("WAAP domain %q not found: create it in the test account or set %s",
		WaapTestDomainName, EnvWaapDomainID)
}
