// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_certificate_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cdn_certificate"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCDNCertificateModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cdn_certificate.CDNCertificateModel)(nil)
	schema := cdn_certificate.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
