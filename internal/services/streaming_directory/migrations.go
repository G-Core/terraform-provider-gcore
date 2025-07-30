// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_directory

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.ResourceWithUpgradeState = (*StreamingDirectoryResource)(nil)

func (r *StreamingDirectoryResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{}
}
