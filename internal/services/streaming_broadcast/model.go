// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_broadcast

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type StreamingBroadcastModel struct {
	BroadcastID        types.Int64                       `tfsdk:"broadcast_id" path:"broadcast_id,optional"`
	Broadcast          *StreamingBroadcastBroadcastModel `tfsdk:"broadcast" json:"broadcast,optional,no_refresh"`
	AdID               types.Int64                       `tfsdk:"ad_id" json:"ad_id,computed"`
	CustomIframeURL    types.String                      `tfsdk:"custom_iframe_url" json:"custom_iframe_url,computed"`
	Name               types.String                      `tfsdk:"name" json:"name,computed"`
	PendingMessage     types.String                      `tfsdk:"pending_message" json:"pending_message,computed"`
	PlayerID           types.Int64                       `tfsdk:"player_id" json:"player_id,computed"`
	Poster             types.String                      `tfsdk:"poster" json:"poster,computed"`
	ShareURL           types.String                      `tfsdk:"share_url" json:"share_url,computed"`
	ShowDvrAfterFinish types.Bool                        `tfsdk:"show_dvr_after_finish" json:"show_dvr_after_finish,computed"`
	Status             types.String                      `tfsdk:"status" json:"status,computed"`
	StreamIDs          customfield.List[types.Int64]     `tfsdk:"stream_ids" json:"stream_ids,computed"`
}

func (m StreamingBroadcastModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m StreamingBroadcastModel) MarshalJSONForUpdate(state StreamingBroadcastModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type StreamingBroadcastBroadcastModel struct {
	Name               types.String   `tfsdk:"name" json:"name,required"`
	AdID               types.Int64    `tfsdk:"ad_id" json:"ad_id,optional"`
	CustomIframeURL    types.String   `tfsdk:"custom_iframe_url" json:"custom_iframe_url,optional"`
	PendingMessage     types.String   `tfsdk:"pending_message" json:"pending_message,optional"`
	PlayerID           types.Int64    `tfsdk:"player_id" json:"player_id,optional"`
	Poster             types.String   `tfsdk:"poster" json:"poster,optional"`
	ShareURL           types.String   `tfsdk:"share_url" json:"share_url,optional"`
	ShowDvrAfterFinish types.Bool     `tfsdk:"show_dvr_after_finish" json:"show_dvr_after_finish,optional"`
	Status             types.String   `tfsdk:"status" json:"status,optional"`
	StreamIDs          *[]types.Int64 `tfsdk:"stream_ids" json:"stream_ids,optional"`
}
