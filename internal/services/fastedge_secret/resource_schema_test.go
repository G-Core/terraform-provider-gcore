// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_secret_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/fastedge_secret"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestFastedgeSecretModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*fastedge_secret.FastedgeSecretModel)(nil)
	schema := fastedge_secret.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
