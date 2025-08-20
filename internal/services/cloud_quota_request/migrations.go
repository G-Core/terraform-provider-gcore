// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_quota_request

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.ResourceWithUpgradeState = (*CloudQuotaRequestResource)(nil)

func (r *CloudQuotaRequestResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{}
}
