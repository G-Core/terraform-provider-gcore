// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_kv_store

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.ResourceWithUpgradeState = (*FastedgeKvStoreResource)(nil)

func (r *FastedgeKvStoreResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{}
}
