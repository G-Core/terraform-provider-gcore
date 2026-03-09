// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_secret_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_secret"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudSecretModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_secret.CloudSecretModel)(nil)
	schema := cloud_secret.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
