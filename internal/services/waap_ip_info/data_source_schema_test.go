// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_ip_info_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/waap_ip_info"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestWaapIPInfoDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*waap_ip_info.WaapIPInfoDataSourceModel)(nil)
	schema := waap_ip_info.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
