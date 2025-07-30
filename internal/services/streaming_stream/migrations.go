// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_stream

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.ResourceWithUpgradeState = (*StreamingStreamResource)(nil)

func (r *StreamingStreamResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{}
}
