// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_instance_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_instance"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudInstanceDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_instance.CloudInstanceDataSourceModel)(nil)
	schema := cloud_instance.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
