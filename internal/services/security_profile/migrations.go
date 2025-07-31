// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package security_profile

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.ResourceWithUpgradeState = (*SecurityProfileResource)(nil)

func (r *SecurityProfileResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{}
}
