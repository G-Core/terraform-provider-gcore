// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_registry_user

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.ResourceWithUpgradeState = (*CloudRegistryUserResource)(nil)

func (r *CloudRegistryUserResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{}
}
