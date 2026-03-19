// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_project_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_project"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudProjectModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_project.CloudProjectModel)(nil)
	schema := cloud_project.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
