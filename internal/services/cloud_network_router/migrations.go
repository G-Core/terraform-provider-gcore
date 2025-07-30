// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_network_router

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.ResourceWithUpgradeState = (*CloudNetworkRouterResource)(nil)

func (r *CloudNetworkRouterResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{}
}
