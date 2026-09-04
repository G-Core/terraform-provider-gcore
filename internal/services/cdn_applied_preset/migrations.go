// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_applied_preset

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.ResourceWithUpgradeState = (*CDNAppliedPresetResource)(nil)

func (r *CDNAppliedPresetResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{}
}
