// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_certificate_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cdn_certificate"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCDNCertificateModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cdn_certificate.CDNCertificateModel)(nil)
	schema := cdn_certificate.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
