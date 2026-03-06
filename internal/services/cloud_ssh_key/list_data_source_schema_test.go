// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_ssh_key_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_ssh_key"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudSSHKeysDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_ssh_key.CloudSSHKeysDataSourceModel)(nil)
	schema := cloud_ssh_key.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
