// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_resource_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cdn_resource"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCDNResourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cdn_resource.CDNResourceModel)(nil)
	schema := cdn_resource.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
