// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_instance_image_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_instance_image"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudInstanceImageModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_instance_image.CloudInstanceImageModel)(nil)
	schema := cloud_instance_image.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
