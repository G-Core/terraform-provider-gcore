// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_directory

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type StreamingDirectoryModel struct {
	ID         types.Int64                                                `tfsdk:"id" json:"id,computed"`
	Name       types.String                                               `tfsdk:"name" json:"name,required,no_refresh"`
	ParentID   types.Int64                                                `tfsdk:"parent_id" json:"parent_id,optional,no_refresh"`
	CreatedAt  types.String                                               `tfsdk:"created_at" json:"created_at,computed,no_refresh"`
	ItemsCount types.Int64                                                `tfsdk:"items_count" json:"items_count,computed,no_refresh"`
	UpdatedAt  types.String                                               `tfsdk:"updated_at" json:"updated_at,computed,no_refresh"`
	Directory  customfield.NestedObject[StreamingDirectoryDirectoryModel] `tfsdk:"directory" json:"directory,computed"`
}

func (m StreamingDirectoryModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m StreamingDirectoryModel) MarshalJSONForUpdate(state StreamingDirectoryModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type StreamingDirectoryDirectoryModel struct {
	ID         types.Int64                                                         `tfsdk:"id" json:"id,computed"`
	CreatedAt  types.String                                                        `tfsdk:"created_at" json:"created_at,computed"`
	ItemsCount types.Int64                                                         `tfsdk:"items_count" json:"items_count,computed"`
	Name       types.String                                                        `tfsdk:"name" json:"name,computed"`
	ParentID   types.Int64                                                         `tfsdk:"parent_id" json:"parent_id,computed"`
	UpdatedAt  types.String                                                        `tfsdk:"updated_at" json:"updated_at,computed"`
	Items      customfield.NestedObjectList[StreamingDirectoryDirectoryItemsModel] `tfsdk:"items" json:"items,computed"`
}

type StreamingDirectoryDirectoryItemsModel struct {
	ID                  types.Int64                                                                        `tfsdk:"id" json:"id,computed"`
	AdID                types.Int64                                                                        `tfsdk:"ad_id" json:"ad_id,computed"`
	CdnViews            types.Int64                                                                        `tfsdk:"cdn_views" json:"cdn_views,computed"`
	ClientID            types.Int64                                                                        `tfsdk:"client_id" json:"client_id,computed"`
	ClientUserID        types.Int64                                                                        `tfsdk:"client_user_id" json:"client_user_id,computed"`
	ConvertedVideos     customfield.NestedObjectList[StreamingDirectoryDirectoryItemsConvertedVideosModel] `tfsdk:"converted_videos" json:"converted_videos,computed"`
	CustomIframeURL     types.String                                                                       `tfsdk:"custom_iframe_url" json:"custom_iframe_url,computed"`
	DashURL             types.String                                                                       `tfsdk:"dash_url" json:"dash_url,computed"`
	Description         types.String                                                                       `tfsdk:"description" json:"description,computed"`
	Duration            types.Int64                                                                        `tfsdk:"duration" json:"duration,computed"`
	Error               types.String                                                                       `tfsdk:"error" json:"error,computed"`
	HlsCmafURL          types.String                                                                       `tfsdk:"hls_cmaf_url" json:"hls_cmaf_url,computed"`
	HlsURL              types.String                                                                       `tfsdk:"hls_url" json:"hls_url,computed"`
	IframeURL           types.String                                                                       `tfsdk:"iframe_url" json:"iframe_url,computed"`
	Name                types.String                                                                       `tfsdk:"name" json:"name,computed"`
	OriginSize          types.Int64                                                                        `tfsdk:"origin_size" json:"origin_size,computed"`
	OriginURL           types.String                                                                       `tfsdk:"origin_url" json:"origin_url,computed"`
	OriginVideoDuration types.Int64                                                                        `tfsdk:"origin_video_duration" json:"origin_video_duration,computed"`
	Poster              types.String                                                                       `tfsdk:"poster" json:"poster,computed"`
	PosterThumb         types.String                                                                       `tfsdk:"poster_thumb" json:"poster_thumb,computed"`
	Projection          types.String                                                                       `tfsdk:"projection" json:"projection,computed"`
	RecordingStartedAt  types.String                                                                       `tfsdk:"recording_started_at" json:"recording_started_at,computed"`
	Screenshot          types.String                                                                       `tfsdk:"screenshot" json:"screenshot,computed"`
	ScreenshotID        types.Int64                                                                        `tfsdk:"screenshot_id" json:"screenshot_id,computed"`
	Screenshots         customfield.List[types.String]                                                     `tfsdk:"screenshots" json:"screenshots,computed"`
	ShareURL            types.String                                                                       `tfsdk:"share_url" json:"share_url,computed"`
	Slug                types.String                                                                       `tfsdk:"slug" json:"slug,computed"`
	Sprite              types.String                                                                       `tfsdk:"sprite" json:"sprite,computed"`
	SpriteVtt           types.String                                                                       `tfsdk:"sprite_vtt" json:"sprite_vtt,computed"`
	Status              types.String                                                                       `tfsdk:"status" json:"status,computed"`
	StreamID            types.Int64                                                                        `tfsdk:"stream_id" json:"stream_id,computed"`
	Views               types.Int64                                                                        `tfsdk:"views" json:"views,computed"`
	ItemType            types.String                                                                       `tfsdk:"item_type" json:"item_type,computed"`
	CreatedAt           types.String                                                                       `tfsdk:"created_at" json:"created_at,computed"`
	ItemsCount          types.Int64                                                                        `tfsdk:"items_count" json:"items_count,computed"`
	ParentID            types.Int64                                                                        `tfsdk:"parent_id" json:"parent_id,computed"`
	UpdatedAt           types.String                                                                       `tfsdk:"updated_at" json:"updated_at,computed"`
}

type StreamingDirectoryDirectoryItemsConvertedVideosModel struct {
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
