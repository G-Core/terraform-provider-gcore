// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_client_config_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cdn_client_config"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCDNClientConfigModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cdn_client_config.CDNClientConfigModel)(nil)
	schema := cdn_client_config.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
