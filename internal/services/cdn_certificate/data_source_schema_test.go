// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_certificate_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cdn_certificate"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCDNCertificateDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cdn_certificate.CDNCertificateDataSourceModel)(nil)
	schema := cdn_certificate.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
