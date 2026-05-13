// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package storage_access_key

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.ResourceWithUpgradeState = (*StorageAccessKeyResource)(nil)

func (r *StorageAccessKeyResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{}
}
