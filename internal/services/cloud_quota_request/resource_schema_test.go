// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_quota_request_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_quota_request"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudQuotaRequestModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_quota_request.CloudQuotaRequestModel)(nil)
	schema := cloud_quota_request.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
