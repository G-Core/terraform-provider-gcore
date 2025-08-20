// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_stream_overlay

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type StreamingStreamOverlayModel struct {
	StreamID  types.Int64                                                   `tfsdk:"stream_id" path:"stream_id,required"`
	OverlayID types.Int64                                                   `tfsdk:"overlay_id" path:"overlay_id,optional"`
	Body      customfield.NestedObjectList[StreamingStreamOverlayBodyModel] `tfsdk:"body" json:"body,computed_optional,no_refresh"`
	Height    types.Int64                                                   `tfsdk:"height" json:"height,optional"`
	URL       types.String                                                  `tfsdk:"url" json:"url,optional"`
	Width     types.Int64                                                   `tfsdk:"width" json:"width,optional"`
	X         types.Int64                                                   `tfsdk:"x" json:"x,optional"`
	Y         types.Int64                                                   `tfsdk:"y" json:"y,optional"`
	Stretch   types.Bool                                                    `tfsdk:"stretch" json:"stretch,computed_optional"`
	CreatedAt types.String                                                  `tfsdk:"created_at" json:"created_at,computed"`
	ID        types.Int64                                                   `tfsdk:"id" json:"id,computed"`
	UpdatedAt types.String                                                  `tfsdk:"updated_at" json:"updated_at,computed"`
}

func (m StreamingStreamOverlayModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m.Body)
}

func (m StreamingStreamOverlayModel) MarshalJSONForUpdate(state StreamingStreamOverlayModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m.Body, state.Body)
}

type StreamingStreamOverlayBodyModel struct {
	URL     types.String `tfsdk:"url" json:"url,required"`
	Height  types.Int64  `tfsdk:"height" json:"height,optional"`
	Stretch types.Bool   `tfsdk:"stretch" json:"stretch,computed_optional"`
	Width   types.Int64  `tfsdk:"width" json:"width,optional"`
	X       types.Int64  `tfsdk:"x" json:"x,optional"`
	Y       types.Int64  `tfsdk:"y" json:"y,optional"`
}
