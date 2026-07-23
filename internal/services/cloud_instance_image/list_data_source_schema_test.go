// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_instance_image_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_instance_image"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudInstanceImagesDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_instance_image.CloudInstanceImagesDataSourceModel)(nil)
	schema := cloud_instance_image.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
