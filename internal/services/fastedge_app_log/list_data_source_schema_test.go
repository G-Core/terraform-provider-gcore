// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_app_log_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/fastedge_app_log"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestFastedgeAppLogsDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*fastedge_app_log.FastedgeAppLogsDataSourceModel)(nil)
	schema := fastedge_app_log.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
