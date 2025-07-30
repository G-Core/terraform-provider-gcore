// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_video

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type StreamingVideoModel struct {
	VideoID                        types.Int64                                                      `tfsdk:"video_id" path:"video_id,optional"`
	Video                          customfield.NestedObject[StreamingVideoVideoModel]               `tfsdk:"video" json:"video,computed_optional,no_refresh"`
	ClientUserID                   types.Int64                                                      `tfsdk:"client_user_id" json:"client_user_id,optional"`
	ClipDurationSeconds            types.Int64                                                      `tfsdk:"clip_duration_seconds" json:"clip_duration_seconds,optional,no_refresh"`
	ClipStartSeconds               types.Int64                                                      `tfsdk:"clip_start_seconds" json:"clip_start_seconds,optional,no_refresh"`
	CustomIframeURL                types.String                                                     `tfsdk:"custom_iframe_url" json:"custom_iframe_url,optional"`
	Description                    types.String                                                     `tfsdk:"description" json:"description,optional"`
	DirectoryID                    types.Int64                                                      `tfsdk:"directory_id" json:"directory_id,optional,no_refresh"`
	Name                           types.String                                                     `tfsdk:"name" json:"name,optional"`
	OriginHTTPHeaders              types.String                                                     `tfsdk:"origin_http_headers" json:"origin_http_headers,optional,no_refresh"`
	OriginURL                      types.String                                                     `tfsdk:"origin_url" json:"origin_url,optional"`
	Poster                         types.String                                                     `tfsdk:"poster" json:"poster,optional"`
	Projection                     types.String                                                     `tfsdk:"projection" json:"projection,optional"`
	QualitySetID                   types.Int64                                                      `tfsdk:"quality_set_id" json:"quality_set_id,optional,no_refresh"`
	RemotePosterURL                types.String                                                     `tfsdk:"remote_poster_url" json:"remote_poster_url,optional,no_refresh"`
	ShareURL                       types.String                                                     `tfsdk:"share_url" json:"share_url,optional"`
	AutoTranscribeAudioLanguage    types.String                                                     `tfsdk:"auto_transcribe_audio_language" json:"auto_transcribe_audio_language,computed_optional,no_refresh"`
	AutoTranslateSubtitlesLanguage types.String                                                     `tfsdk:"auto_translate_subtitles_language" json:"auto_translate_subtitles_language,computed_optional,no_refresh"`
	Priority                       types.Int64                                                      `tfsdk:"priority" json:"priority,computed_optional,no_refresh"`
	RemovePoster                   types.Bool                                                       `tfsdk:"remove_poster" json:"remove_poster,computed_optional,no_refresh"`
	ScreenshotID                   types.Int64                                                      `tfsdk:"screenshot_id" json:"screenshot_id,computed_optional"`
	SourceBitrateLimit             types.Bool                                                       `tfsdk:"source_bitrate_limit" json:"source_bitrate_limit,computed_optional,no_refresh"`
	AdID                           types.Int64                                                      `tfsdk:"ad_id" json:"ad_id,computed"`
	CdnViews                       types.Int64                                                      `tfsdk:"cdn_views" json:"cdn_views,computed"`
	ClientID                       types.Int64                                                      `tfsdk:"client_id" json:"client_id,computed"`
	DashURL                        types.String                                                     `tfsdk:"dash_url" json:"dash_url,computed"`
	Duration                       types.Int64                                                      `tfsdk:"duration" json:"duration,computed"`
	Error                          types.String                                                     `tfsdk:"error" json:"error,computed"`
	HlsCmafURL                     types.String                                                     `tfsdk:"hls_cmaf_url" json:"hls_cmaf_url,computed"`
	HlsURL                         types.String                                                     `tfsdk:"hls_url" json:"hls_url,computed"`
	ID                             types.Int64                                                      `tfsdk:"id" json:"id,computed"`
	IframeURL                      types.String                                                     `tfsdk:"iframe_url" json:"iframe_url,computed"`
	OriginSize                     types.Int64                                                      `tfsdk:"origin_size" json:"origin_size,computed"`
	OriginVideoDuration            types.Int64                                                      `tfsdk:"origin_video_duration" json:"origin_video_duration,computed"`
	PosterThumb                    types.String                                                     `tfsdk:"poster_thumb" json:"poster_thumb,computed"`
	RecordingStartedAt             types.String                                                     `tfsdk:"recording_started_at" json:"recording_started_at,computed"`
	Screenshot                     types.String                                                     `tfsdk:"screenshot" json:"screenshot,computed"`
	Slug                           types.String                                                     `tfsdk:"slug" json:"slug,computed"`
	Sprite                         types.String                                                     `tfsdk:"sprite" json:"sprite,computed"`
	SpriteVtt                      types.String                                                     `tfsdk:"sprite_vtt" json:"sprite_vtt,computed"`
	Status                         types.String                                                     `tfsdk:"status" json:"status,computed"`
	StreamID                       types.Int64                                                      `tfsdk:"stream_id" json:"stream_id,computed"`
	Views                          types.Int64                                                      `tfsdk:"views" json:"views,computed"`
	Screenshots                    customfield.List[types.String]                                   `tfsdk:"screenshots" json:"screenshots,computed"`
	ConvertedVideos                customfield.NestedObjectList[StreamingVideoConvertedVideosModel] `tfsdk:"converted_videos" json:"converted_videos,computed"`
}

func (m StreamingVideoModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m StreamingVideoModel) MarshalJSONForUpdate(state StreamingVideoModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type StreamingVideoVideoModel struct {
	Name                           types.String `tfsdk:"name" json:"name,required"`
	AutoTranscribeAudioLanguage    types.String `tfsdk:"auto_transcribe_audio_language" json:"auto_transcribe_audio_language,computed_optional"`
	AutoTranslateSubtitlesLanguage types.String `tfsdk:"auto_translate_subtitles_language" json:"auto_translate_subtitles_language,computed_optional"`
	ClientUserID                   types.Int64  `tfsdk:"client_user_id" json:"client_user_id,optional"`
	ClipDurationSeconds            types.Int64  `tfsdk:"clip_duration_seconds" json:"clip_duration_seconds,optional"`
	ClipStartSeconds               types.Int64  `tfsdk:"clip_start_seconds" json:"clip_start_seconds,optional"`
	CustomIframeURL                types.String `tfsdk:"custom_iframe_url" json:"custom_iframe_url,optional"`
	Description                    types.String `tfsdk:"description" json:"description,optional"`
	DirectoryID                    types.Int64  `tfsdk:"directory_id" json:"directory_id,optional"`
	OriginHTTPHeaders              types.String `tfsdk:"origin_http_headers" json:"origin_http_headers,optional"`
	OriginURL                      types.String `tfsdk:"origin_url" json:"origin_url,optional"`
	Poster                         types.String `tfsdk:"poster" json:"poster,optional"`
	Priority                       types.Int64  `tfsdk:"priority" json:"priority,computed_optional"`
	Projection                     types.String `tfsdk:"projection" json:"projection,optional"`
	QualitySetID                   types.Int64  `tfsdk:"quality_set_id" json:"quality_set_id,optional"`
	RemotePosterURL                types.String `tfsdk:"remote_poster_url" json:"remote_poster_url,optional"`
	RemovePoster                   types.Bool   `tfsdk:"remove_poster" json:"remove_poster,computed_optional"`
	ScreenshotID                   types.Int64  `tfsdk:"screenshot_id" json:"screenshot_id,computed_optional"`
	ShareURL                       types.String `tfsdk:"share_url" json:"share_url,optional"`
	SourceBitrateLimit             types.Bool   `tfsdk:"source_bitrate_limit" json:"source_bitrate_limit,computed_optional"`
}

type StreamingVideoConvertedVideosModel struct {
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
