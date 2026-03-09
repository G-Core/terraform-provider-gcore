// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_ssh_key_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_ssh_key"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudSSHKeyModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_ssh_key.CloudSSHKeyModel)(nil)
	schema := cloud_ssh_key.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
