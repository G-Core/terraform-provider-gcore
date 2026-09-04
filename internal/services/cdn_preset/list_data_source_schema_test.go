// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_preset_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cdn_preset"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCDNPresetsDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cdn_preset.CDNPresetsDataSourceModel)(nil)
	schema := cdn_preset.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
