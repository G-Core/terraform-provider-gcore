// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_volume_snapshot_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_volume_snapshot"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudVolumeSnapshotsDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_volume_snapshot.CloudVolumeSnapshotsDataSourceModel)(nil)
	schema := cloud_volume_snapshot.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
