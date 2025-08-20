// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package security_event_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/security_event"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestSecurityEventsDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*security_event.SecurityEventsDataSourceModel)(nil)
	schema := security_event.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
