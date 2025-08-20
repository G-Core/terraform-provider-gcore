// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_volume_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_volume"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudVolumesDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_volume.CloudVolumesDataSourceModel)(nil)
	schema := cloud_volume.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
