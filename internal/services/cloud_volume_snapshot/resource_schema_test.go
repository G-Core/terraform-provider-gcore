// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_volume_snapshot_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_volume_snapshot"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudVolumeSnapshotModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_volume_snapshot.CloudVolumeSnapshotModel)(nil)
	schema := cloud_volume_snapshot.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
