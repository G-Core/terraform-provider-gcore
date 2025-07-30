// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_audit_log_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_audit_log"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudAuditLogsDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_audit_log.CloudAuditLogsDataSourceModel)(nil)
	schema := cloud_audit_log.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
