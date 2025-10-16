// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer_pool_member

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.ResourceWithUpgradeState = (*CloudLoadBalancerPoolMemberResource)(nil)

func (r *CloudLoadBalancerPoolMemberResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{}
}
