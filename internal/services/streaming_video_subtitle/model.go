// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_video_subtitle

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
)

type StreamingVideoSubtitleModel struct {
	VideoID  types.Int64          `tfsdk:"video_id" path:"video_id,required"`
	ID       types.Int64          `tfsdk:"id" path:"id,optional"`
	Body     jsontypes.Normalized `tfsdk:"body" json:"body,required,no_refresh"`
	Language types.String         `tfsdk:"language" json:"language,optional"`
	Name     types.String         `tfsdk:"name" json:"name,optional"`
	Vtt      types.String         `tfsdk:"vtt" json:"vtt,optional"`
}

func (m StreamingVideoSubtitleModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m.Body)
}

func (m StreamingVideoSubtitleModel) MarshalJSONForUpdate(state StreamingVideoSubtitleModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m.Body, state.Body)
}
