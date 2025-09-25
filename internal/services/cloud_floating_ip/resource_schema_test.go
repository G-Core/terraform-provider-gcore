// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_floating_ip_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_floating_ip"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudFloatingIPModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_floating_ip.CloudFloatingIPModel)(nil)
	schema := cloud_floating_ip.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
