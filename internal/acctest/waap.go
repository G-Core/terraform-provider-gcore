package acctest

import (
	"context"
	"os"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/sweep"
)

// WaapDomainID returns the WAAP domain the WAAP rule tests attach rules to; see
// sweep.WaapDomainID. The sweepers resolve the same domain, so whatever a test
// leaks is where the sweep looks.
func WaapDomainID(t *testing.T) int64 {
	t.Helper()

	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	client, err := NewGcoreClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	id, err := sweep.WaapDomainID(context.Background(), client)
	if err != nil {
		t.Fatal(err)
	}
	return id
}
