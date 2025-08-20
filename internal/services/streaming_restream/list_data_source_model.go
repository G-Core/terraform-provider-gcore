// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_restream

import (
	"context"

	"github.com/G-Core/gcore-go/streaming"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type StreamingRestreamsItemsListDataSourceEnvelope struct {
	Items customfield.NestedObjectList[StreamingRestreamsItemsDataSourceModel] `json:"items,computed"`
}

type StreamingRestreamsDataSourceModel struct {
	MaxItems types.Int64                                                          `tfsdk:"max_items"`
	Items    customfield.NestedObjectList[StreamingRestreamsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *StreamingRestreamsDataSourceModel) toListParams(_ context.Context) (params streaming.RestreamListParams, diags diag.Diagnostics) {
	params = streaming.RestreamListParams{}

	return
}

type StreamingRestreamsItemsDataSourceModel struct {
	Active       types.Bool   `tfsdk:"active" json:"active,computed"`
	ClientUserID types.Int64  `tfsdk:"client_user_id" json:"client_user_id,computed"`
	Live         types.Bool   `tfsdk:"live" json:"live,computed"`
	Name         types.String `tfsdk:"name" json:"name,computed"`
	StreamID     types.Int64  `tfsdk:"stream_id" json:"stream_id,computed"`
	Uri          types.String `tfsdk:"uri" json:"uri,computed"`
}
