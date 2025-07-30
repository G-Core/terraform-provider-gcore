// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_api_path_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/waap_domain_api_path"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestWaapDomainAPIPathModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*waap_domain_api_path.WaapDomainAPIPathModel)(nil)
	schema := waap_domain_api_path.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
