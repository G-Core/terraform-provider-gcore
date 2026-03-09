// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_app_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/fastedge_app"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestFastedgeAppModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*fastedge_app.FastedgeAppModel)(nil)
	schema := fastedge_app.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
