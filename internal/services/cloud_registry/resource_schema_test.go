// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_registry_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_registry"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudRegistryModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_registry.CloudRegistryModel)(nil)
	schema := cloud_registry.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
