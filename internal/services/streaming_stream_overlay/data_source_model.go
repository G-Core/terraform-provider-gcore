// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_stream_overlay

import (
	"context"

	"github.com/G-Core/gcore-go/streaming"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type StreamingStreamOverlayDataSourceModel struct {
	OverlayID types.Int64  `tfsdk:"overlay_id" path:"overlay_id,required"`
	StreamID  types.Int64  `tfsdk:"stream_id" path:"stream_id,required"`
	CreatedAt types.String `tfsdk:"created_at" json:"created_at,computed"`
	Height    types.Int64  `tfsdk:"height" json:"height,computed"`
	ID        types.Int64  `tfsdk:"id" json:"id,computed"`
	Stretch   types.Bool   `tfsdk:"stretch" json:"stretch,computed"`
	UpdatedAt types.String `tfsdk:"updated_at" json:"updated_at,computed"`
	URL       types.String `tfsdk:"url" json:"url,computed"`
	Width     types.Int64  `tfsdk:"width" json:"width,computed"`
	X         types.Int64  `tfsdk:"x" json:"x,computed"`
	Y         types.Int64  `tfsdk:"y" json:"y,computed"`
}

func (m *StreamingStreamOverlayDataSourceModel) toReadParams(_ context.Context) (params streaming.StreamOverlayGetParams, diags diag.Diagnostics) {
	params = streaming.StreamOverlayGetParams{
		StreamID: m.StreamID.ValueInt64(),
	}

	return
}
