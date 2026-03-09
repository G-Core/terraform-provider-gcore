// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_security_group_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_security_group"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudSecurityGroupModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_security_group.CloudSecurityGroupModel)(nil)
	schema := cloud_security_group.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
