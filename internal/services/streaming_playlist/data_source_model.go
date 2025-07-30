// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_playlist

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type StreamingPlaylistDataSourceModel struct {
	PlaylistID   types.Int64                   `tfsdk:"playlist_id" path:"playlist_id,required"`
	Active       types.Bool                    `tfsdk:"active" json:"active,computed"`
	AdID         types.Int64                   `tfsdk:"ad_id" json:"ad_id,computed"`
	ClientID     types.Int64                   `tfsdk:"client_id" json:"client_id,computed"`
	ClientUserID types.Int64                   `tfsdk:"client_user_id" json:"client_user_id,computed"`
	Countdown    types.Bool                    `tfsdk:"countdown" json:"countdown,computed"`
	HlsCmafURL   types.String                  `tfsdk:"hls_cmaf_url" json:"hls_cmaf_url,computed"`
	HlsURL       types.String                  `tfsdk:"hls_url" json:"hls_url,computed"`
	IframeURL    types.String                  `tfsdk:"iframe_url" json:"iframe_url,computed"`
	Loop         types.Bool                    `tfsdk:"loop" json:"loop,computed"`
	Name         types.String                  `tfsdk:"name" json:"name,computed"`
	PlayerID     types.Int64                   `tfsdk:"player_id" json:"player_id,computed"`
	PlaylistType types.String                  `tfsdk:"playlist_type" json:"playlist_type,computed"`
	StartTime    types.String                  `tfsdk:"start_time" json:"start_time,computed"`
	VideoIDs     customfield.List[types.Int64] `tfsdk:"video_ids" json:"video_ids,computed"`
}
