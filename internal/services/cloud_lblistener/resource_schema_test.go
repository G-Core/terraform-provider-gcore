// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_lblistener_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_lblistener"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudLblistenerModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_lblistener.CloudLblistenerModel)(nil)
	schema := cloud_lblistener.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
