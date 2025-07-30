// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_registry_user_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_registry_user"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudRegistryUserModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_registry_user.CloudRegistryUserModel)(nil)
	schema := cloud_registry_user.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
