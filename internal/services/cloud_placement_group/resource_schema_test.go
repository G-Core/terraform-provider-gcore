// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_placement_group_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_placement_group"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudPlacementGroupModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_placement_group.CloudPlacementGroupModel)(nil)
	schema := cloud_placement_group.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
