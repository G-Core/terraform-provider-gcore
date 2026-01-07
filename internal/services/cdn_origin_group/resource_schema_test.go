// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_origin_group_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cdn_origin_group"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCdnOriginGroupModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cdn_origin_group.CdnOriginGroupModel)(nil)
	schema := cdn_origin_group.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
