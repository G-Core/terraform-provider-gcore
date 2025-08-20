// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_setting_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/waap_domain_setting"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestWaapDomainSettingDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*waap_domain_setting.WaapDomainSettingDataSourceModel)(nil)
	schema := waap_domain_setting.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
