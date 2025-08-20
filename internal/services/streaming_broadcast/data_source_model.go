// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_broadcast

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type StreamingBroadcastDataSourceModel struct {
	BroadcastID        types.Int64                   `tfsdk:"broadcast_id" path:"broadcast_id,required"`
	AdID               types.Int64                   `tfsdk:"ad_id" json:"ad_id,computed"`
	CustomIframeURL    types.String                  `tfsdk:"custom_iframe_url" json:"custom_iframe_url,computed"`
	Name               types.String                  `tfsdk:"name" json:"name,computed"`
	PendingMessage     types.String                  `tfsdk:"pending_message" json:"pending_message,computed"`
	PlayerID           types.Int64                   `tfsdk:"player_id" json:"player_id,computed"`
	Poster             types.String                  `tfsdk:"poster" json:"poster,computed"`
	ShareURL           types.String                  `tfsdk:"share_url" json:"share_url,computed"`
	ShowDvrAfterFinish types.Bool                    `tfsdk:"show_dvr_after_finish" json:"show_dvr_after_finish,computed"`
	Status             types.String                  `tfsdk:"status" json:"status,computed"`
	StreamIDs          customfield.List[types.Int64] `tfsdk:"stream_ids" json:"stream_ids,computed"`
}
