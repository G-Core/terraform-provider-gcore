// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_lbpool_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_lbpool"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudLbpoolModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_lbpool.CloudLbpoolModel)(nil)
	schema := cloud_lbpool.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
