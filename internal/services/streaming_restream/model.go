// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_restream

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
)

type StreamingRestreamModel struct {
	RestreamID   types.Int64                     `tfsdk:"restream_id" path:"restream_id,optional"`
	Restream     *StreamingRestreamRestreamModel `tfsdk:"restream" json:"restream,optional,no_refresh"`
	Active       types.Bool                      `tfsdk:"active" json:"active,computed"`
	ClientUserID types.Int64                     `tfsdk:"client_user_id" json:"client_user_id,computed"`
	Live         types.Bool                      `tfsdk:"live" json:"live,computed"`
	Name         types.String                    `tfsdk:"name" json:"name,computed"`
	StreamID     types.Int64                     `tfsdk:"stream_id" json:"stream_id,computed"`
	Uri          types.String                    `tfsdk:"uri" json:"uri,computed"`
}

func (m StreamingRestreamModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m StreamingRestreamModel) MarshalJSONForUpdate(state StreamingRestreamModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type StreamingRestreamRestreamModel struct {
	Active       types.Bool   `tfsdk:"active" json:"active,optional"`
	ClientUserID types.Int64  `tfsdk:"client_user_id" json:"client_user_id,optional"`
	Live         types.Bool   `tfsdk:"live" json:"live,optional"`
	Name         types.String `tfsdk:"name" json:"name,optional"`
	StreamID     types.Int64  `tfsdk:"stream_id" json:"stream_id,optional"`
	Uri          types.String `tfsdk:"uri" json:"uri,optional"`
}
