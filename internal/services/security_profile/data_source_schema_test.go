// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package security_profile_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/security_profile"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestSecurityProfileDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*security_profile.SecurityProfileDataSourceModel)(nil)
	schema := security_profile.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
