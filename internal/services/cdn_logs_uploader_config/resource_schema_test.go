// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_logs_uploader_config_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cdn_logs_uploader_config"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCDNLogsUploaderConfigModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cdn_logs_uploader_config.CDNLogsUploaderConfigModel)(nil)
	schema := cdn_logs_uploader_config.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
