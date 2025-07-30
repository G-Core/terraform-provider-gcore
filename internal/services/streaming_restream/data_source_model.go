// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_restream

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type StreamingRestreamDataSourceModel struct {
	RestreamID   types.Int64  `tfsdk:"restream_id" path:"restream_id,required"`
	Active       types.Bool   `tfsdk:"active" json:"active,computed"`
	ClientUserID types.Int64  `tfsdk:"client_user_id" json:"client_user_id,computed"`
	Live         types.Bool   `tfsdk:"live" json:"live,computed"`
	Name         types.String `tfsdk:"name" json:"name,computed"`
	StreamID     types.Int64  `tfsdk:"stream_id" json:"stream_id,computed"`
	Uri          types.String `tfsdk:"uri" json:"uri,computed"`
}
