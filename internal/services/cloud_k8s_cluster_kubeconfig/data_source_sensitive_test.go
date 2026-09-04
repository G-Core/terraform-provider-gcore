package cloud_k8s_cluster_kubeconfig_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_k8s_cluster_kubeconfig"
)

// TestDataSourceSensitiveAttributes pins the sensitivity of every attribute of
// the kubeconfig data source.
//
// The kubeconfig and its credential components must never appear in plan,
// apply, or CI output. Sensitivity is currently applied as custom code on the
// generated schema; the source API spec can drive it through
// x-stainless-sensitive instead, at which point this test guards the
// regenerated schema unchanged.
//
// host, created_at, and expires_at are deliberately not sensitive: they are
// connection and lifetime metadata used for diagnostics, not credentials.
func TestDataSourceSensitiveAttributes(t *testing.T) {
	t.Parallel()

	want := map[string]bool{
		"client_certificate":     true,
		"client_key":             true,
		"cluster_ca_certificate": true,
		"config":                 true,
		"cluster_name":           false,
		"created_at":             false,
		"expires_at":             false,
		"host":                   false,
		"project_id":             false,
		"region_id":              false,
	}

	attrs := cloud_k8s_cluster_kubeconfig.DataSourceSchema(context.TODO()).Attributes
	if len(attrs) != len(want) {
		t.Errorf("schema has %d attributes, expectations cover %d: keep this test in sync with the schema", len(attrs), len(want))
	}

	for name, wantSensitive := range want {
		attr, ok := attrs[name]
		if !ok {
			t.Errorf("attribute %q missing from schema", name)
			continue
		}
		if got := attr.IsSensitive(); got != wantSensitive {
			t.Errorf("attribute %q: IsSensitive() = %t, want %t", name, got, wantSensitive)
		}
	}
}
