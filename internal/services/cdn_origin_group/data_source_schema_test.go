// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_origin_group_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cdn_origin_group"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCDNOriginGroupDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cdn_origin_group.CDNOriginGroupDataSourceModel)(nil)
	schema := cdn_origin_group.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
