// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_playlist

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
)

type StreamingPlaylistModel struct {
	PlaylistID   types.Int64    `tfsdk:"playlist_id" path:"playlist_id,optional"`
	AdID         types.Int64    `tfsdk:"ad_id" json:"ad_id,optional"`
	ClientID     types.Int64    `tfsdk:"client_id" json:"client_id,optional"`
	ClientUserID types.Int64    `tfsdk:"client_user_id" json:"client_user_id,optional"`
	Countdown    types.Bool     `tfsdk:"countdown" json:"countdown,optional"`
	HlsCmafURL   types.String   `tfsdk:"hls_cmaf_url" json:"hls_cmaf_url,optional"`
	HlsURL       types.String   `tfsdk:"hls_url" json:"hls_url,optional"`
	IframeURL    types.String   `tfsdk:"iframe_url" json:"iframe_url,optional"`
	Name         types.String   `tfsdk:"name" json:"name,optional"`
	PlayerID     types.Int64    `tfsdk:"player_id" json:"player_id,optional"`
	PlaylistType types.String   `tfsdk:"playlist_type" json:"playlist_type,optional"`
	StartTime    types.String   `tfsdk:"start_time" json:"start_time,optional"`
	VideoIDs     *[]types.Int64 `tfsdk:"video_ids" json:"video_ids,optional"`
	Active       types.Bool     `tfsdk:"active" json:"active,computed_optional"`
	Loop         types.Bool     `tfsdk:"loop" json:"loop,computed_optional"`
	ID           types.Int64    `tfsdk:"id" json:"id,computed,no_refresh"`
}

func (m StreamingPlaylistModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m StreamingPlaylistModel) MarshalJSONForUpdate(state StreamingPlaylistModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}
