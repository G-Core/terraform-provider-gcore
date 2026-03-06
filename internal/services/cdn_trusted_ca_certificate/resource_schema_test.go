// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_trusted_ca_certificate_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cdn_trusted_ca_certificate"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCDNTrustedCaCertificateModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cdn_trusted_ca_certificate.CDNTrustedCaCertificateModel)(nil)
	schema := cdn_trusted_ca_certificate.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
