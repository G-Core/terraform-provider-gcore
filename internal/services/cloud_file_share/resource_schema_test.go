// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_file_share_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_file_share"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudFileShareModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_file_share.CloudFileShareModel)(nil)
	schema := cloud_file_share.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
