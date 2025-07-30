// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_video

import (
	"context"

	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/gcore-go/streaming"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type StreamingVideosItemsListDataSourceEnvelope struct {
	Items customfield.NestedObjectList[StreamingVideosItemsDataSourceModel] `json:"items,computed"`
}

type StreamingVideosDataSourceModel struct {
	ClientUserID types.Int64                                                       `tfsdk:"client_user_id" query:"client_user_id,optional"`
	Fields       types.String                                                      `tfsdk:"fields" query:"fields,optional"`
	ID           types.String                                                      `tfsdk:"id" query:"id,optional"`
	PerPage      types.Int64                                                       `tfsdk:"per_page" query:"per_page,optional"`
	Search       types.String                                                      `tfsdk:"search" query:"search,optional"`
	Status       types.String                                                      `tfsdk:"status" query:"status,optional"`
	StreamID     types.Int64                                                       `tfsdk:"stream_id" query:"stream_id,optional"`
	MaxItems     types.Int64                                                       `tfsdk:"max_items"`
	Items        customfield.NestedObjectList[StreamingVideosItemsDataSourceModel] `tfsdk:"items"`
}

func (m *StreamingVideosDataSourceModel) toListParams(_ context.Context) (params streaming.VideoListParams, diags diag.Diagnostics) {
	params = streaming.VideoListParams{}

	if !m.ID.IsNull() {
		params.ID = param.NewOpt(m.ID.ValueString())
	}
	if !m.ClientUserID.IsNull() {
		params.ClientUserID = param.NewOpt(m.ClientUserID.ValueInt64())
	}
	if !m.Fields.IsNull() {
		params.Fields = param.NewOpt(m.Fields.ValueString())
	}
	if !m.PerPage.IsNull() {
		params.PerPage = param.NewOpt(m.PerPage.ValueInt64())
	}
	if !m.Search.IsNull() {
		params.Search = param.NewOpt(m.Search.ValueString())
	}
	if !m.Status.IsNull() {
		params.Status = param.NewOpt(m.Status.ValueString())
	}
	if !m.StreamID.IsNull() {
		params.StreamID = param.NewOpt(m.StreamID.ValueInt64())
	}

	return
}

type StreamingVideosItemsDataSourceModel struct {
	ID                  types.Int64                                                                 `tfsdk:"id" json:"id,computed"`
	AdID                types.Int64                                                                 `tfsdk:"ad_id" json:"ad_id,computed"`
	CdnViews            types.Int64                                                                 `tfsdk:"cdn_views" json:"cdn_views,computed"`
	ClientID            types.Int64                                                                 `tfsdk:"client_id" json:"client_id,computed"`
	ClientUserID        types.Int64                                                                 `tfsdk:"client_user_id" json:"client_user_id,computed"`
	ConvertedVideos     customfield.NestedObjectList[StreamingVideosConvertedVideosDataSourceModel] `tfsdk:"converted_videos" json:"converted_videos,computed"`
	CustomIframeURL     types.String                                                                `tfsdk:"custom_iframe_url" json:"custom_iframe_url,computed"`
	DashURL             types.String                                                                `tfsdk:"dash_url" json:"dash_url,computed"`
	Description         types.String                                                                `tfsdk:"description" json:"description,computed"`
	Duration            types.Int64                                                                 `tfsdk:"duration" json:"duration,computed"`
	Error               types.String                                                                `tfsdk:"error" json:"error,computed"`
	HlsCmafURL          types.String                                                                `tfsdk:"hls_cmaf_url" json:"hls_cmaf_url,computed"`
	HlsURL              types.String                                                                `tfsdk:"hls_url" json:"hls_url,computed"`
	IframeURL           types.String                                                                `tfsdk:"iframe_url" json:"iframe_url,computed"`
	Name                types.String                                                                `tfsdk:"name" json:"name,computed"`
	OriginSize          types.Int64                                                                 `tfsdk:"origin_size" json:"origin_size,computed"`
	OriginURL           types.String                                                                `tfsdk:"origin_url" json:"origin_url,computed"`
	OriginVideoDuration types.Int64                                                                 `tfsdk:"origin_video_duration" json:"origin_video_duration,computed"`
	Poster              types.String                                                                `tfsdk:"poster" json:"poster,computed"`
	PosterThumb         types.String                                                                `tfsdk:"poster_thumb" json:"poster_thumb,computed"`
	Projection          types.String                                                                `tfsdk:"projection" json:"projection,computed"`
	RecordingStartedAt  types.String                                                                `tfsdk:"recording_started_at" json:"recording_started_at,computed"`
	Screenshot          types.String                                                                `tfsdk:"screenshot" json:"screenshot,computed"`
	ScreenshotID        types.Int64                                                                 `tfsdk:"screenshot_id" json:"screenshot_id,computed"`
	Screenshots         customfield.List[types.String]                                              `tfsdk:"screenshots" json:"screenshots,computed"`
	ShareURL            types.String                                                                `tfsdk:"share_url" json:"share_url,computed"`
	Slug                types.String                                                                `tfsdk:"slug" json:"slug,computed"`
	Sprite              types.String                                                                `tfsdk:"sprite" json:"sprite,computed"`
	SpriteVtt           types.String                                                                `tfsdk:"sprite_vtt" json:"sprite_vtt,computed"`
	Status              types.String                                                                `tfsdk:"status" json:"status,computed"`
	StreamID            types.Int64                                                                 `tfsdk:"stream_id" json:"stream_id,computed"`
	Views               types.Int64                                                                 `tfsdk:"views" json:"views,computed"`
}

type StreamingVideosConvertedVideosDataSourceModel struct {
	ID       types.Int64  `tfsdk:"id" json:"id,computed"`
	Error    types.String `tfsdk:"error" json:"error,computed"`
	Height   types.Int64  `tfsdk:"height" json:"height,computed"`
	MP4URL   types.String `tfsdk:"mp4_url" json:"mp4_url,computed"`
	Name     types.String `tfsdk:"name" json:"name,computed"`
	Progress types.Int64  `tfsdk:"progress" json:"progress,computed"`
	Size     types.Int64  `tfsdk:"size" json:"size,computed"`
	Status   types.String `tfsdk:"status" json:"status,computed"`
	Width    types.Int64  `tfsdk:"width" json:"width,computed"`
}
