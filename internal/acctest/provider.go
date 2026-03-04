package acctest

import (
	"os"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// ProtoV6ProviderFactories is a static map used for basic tests
// Following the pattern from terraform-provider-aws
var ProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"gcore": providerserver.NewProtocol6WithError(internal.NewProvider("test")()),
}

// PreCheck validates required environment variables are set for acceptance tests
// This follows the standard HashiCorp pattern used across official providers
func PreCheck(t *testing.T) {
	t.Helper()

	if v := os.Getenv("GCORE_API_KEY"); v == "" {
		t.Fatal("GCORE_API_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("GCORE_CLOUD_PROJECT_ID"); v == "" {
		t.Fatal("GCORE_CLOUD_PROJECT_ID must be set for acceptance tests")
	}
	if v := os.Getenv("GCORE_CLOUD_REGION_ID"); v == "" {
		t.Fatal("GCORE_CLOUD_REGION_ID must be set for acceptance tests")
	}
}
