// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_video_subtitle

import (
	"context"

	"github.com/G-Core/gcore-go/streaming"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type StreamingVideoSubtitleDataSourceModel struct {
	ID       types.Int64  `tfsdk:"id" path:"id,required"`
	VideoID  types.Int64  `tfsdk:"video_id" path:"video_id,required"`
	Language types.String `tfsdk:"language" json:"language,computed"`
	Name     types.String `tfsdk:"name" json:"name,computed"`
	Vtt      types.String `tfsdk:"vtt" json:"vtt,computed"`
}

func (m *StreamingVideoSubtitleDataSourceModel) toReadParams(_ context.Context) (params streaming.VideoSubtitleGetParams, diags diag.Diagnostics) {
	params = streaming.VideoSubtitleGetParams{
		VideoID: m.VideoID.ValueInt64(),
	}

	return
}
