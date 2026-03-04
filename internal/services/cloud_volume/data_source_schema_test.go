// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_volume_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_volume"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudVolumeDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_volume.CloudVolumeDataSourceModel)(nil)
	schema := cloud_volume.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
