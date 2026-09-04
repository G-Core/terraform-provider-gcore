// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_applied_preset_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cdn_applied_preset"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCDNAppliedPresetModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cdn_applied_preset.CDNAppliedPresetModel)(nil)
	schema := cdn_applied_preset.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
