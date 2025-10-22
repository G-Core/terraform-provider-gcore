// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_secret

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.ResourceWithUpgradeState = (*CloudSecretResource)(nil)

func (r *CloudSecretResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{}
}
