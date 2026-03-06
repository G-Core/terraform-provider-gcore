// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_trusted_ca_certificate_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cdn_trusted_ca_certificate"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCDNTrustedCaCertificateDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cdn_trusted_ca_certificate.CDNTrustedCaCertificateDataSourceModel)(nil)
	schema := cdn_trusted_ca_certificate.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
