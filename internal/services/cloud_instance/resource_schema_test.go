// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_instance_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_instance"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudInstanceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_instance.CloudInstanceModel)(nil)
	schema := cloud_instance.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	// project_id and region_id are URL path parameters (not JSON body fields) that use
	// provider defaults when not specified. The schema marks them as computed_optional,
	// but the model uses path tags which don't support the computed_optional semantic.
	errs.Ignore(t, ".@CloudInstanceModel.project_id", ".@CloudInstanceModel.region_id")
	errs.Report(t)
}
