// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_volume_snapshot

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.ResourceWithUpgradeState = (*CloudVolumeSnapshotResource)(nil)

func (r *CloudVolumeSnapshotResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{}
}
