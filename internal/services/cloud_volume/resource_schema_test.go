// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_volume_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_volume"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudVolumeModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_volume.CloudVolumeModel)(nil)
	schema := cloud_volume.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
