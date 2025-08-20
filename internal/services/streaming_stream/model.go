// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_stream

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type StreamingStreamModel struct {
	ID                  types.Int64                                                    `tfsdk:"id" json:"id,computed"`
	Name                types.String                                                   `tfsdk:"name" json:"name,required"`
	CdnID               types.Int64                                                    `tfsdk:"cdn_id" json:"cdn_id,optional"`
	ClientEntityData    types.String                                                   `tfsdk:"client_entity_data" json:"client_entity_data,optional"`
	ClientUserID        types.Int64                                                    `tfsdk:"client_user_id" json:"client_user_id,optional"`
	QualitySetID        types.Int64                                                    `tfsdk:"quality_set_id" json:"quality_set_id,optional"`
	Uri                 types.String                                                   `tfsdk:"uri" json:"uri,optional"`
	Active              types.Bool                                                     `tfsdk:"active" json:"active,computed_optional"`
	AutoRecord          types.Bool                                                     `tfsdk:"auto_record" json:"auto_record,computed_optional"`
	DvrDuration         types.Int64                                                    `tfsdk:"dvr_duration" json:"dvr_duration,computed_optional"`
	DvrEnabled          types.Bool                                                     `tfsdk:"dvr_enabled" json:"dvr_enabled,computed_optional"`
	HlsMpegtsEndlistTag types.Bool                                                     `tfsdk:"hls_mpegts_endlist_tag" json:"hls_mpegts_endlist_tag,computed_optional"`
	HTMLOverlay         types.Bool                                                     `tfsdk:"html_overlay" json:"html_overlay,computed_optional"`
	LowLatencyEnabled   types.Bool                                                     `tfsdk:"low_latency_enabled" json:"low_latency_enabled,computed_optional"`
	Projection          types.String                                                   `tfsdk:"projection" json:"projection,computed_optional"`
	Pull                types.Bool                                                     `tfsdk:"pull" json:"pull,computed_optional"`
	RecordType          types.String                                                   `tfsdk:"record_type" json:"record_type,computed_optional"`
	BroadcastIDs        customfield.List[types.Int64]                                  `tfsdk:"broadcast_ids" json:"broadcast_ids,computed_optional"`
	Stream              customfield.NestedObject[StreamingStreamStreamModel]           `tfsdk:"stream" json:"stream,computed_optional,no_refresh"`
	BackupLive          types.Bool                                                     `tfsdk:"backup_live" json:"backup_live,computed"`
	BackupPushURL       types.String                                                   `tfsdk:"backup_push_url" json:"backup_push_url,computed"`
	BackupPushURLSrt    types.String                                                   `tfsdk:"backup_push_url_srt" json:"backup_push_url_srt,computed"`
	CreatedAt           types.String                                                   `tfsdk:"created_at" json:"created_at,computed"`
	DashURL             types.String                                                   `tfsdk:"dash_url" json:"dash_url,computed"`
	FinishedAtPrimary   types.String                                                   `tfsdk:"finished_at_primary" json:"finished_at_primary,computed"`
	FrameRate           types.Float64                                                  `tfsdk:"frame_rate" json:"frame_rate,computed"`
	HlsCmafURL          types.String                                                   `tfsdk:"hls_cmaf_url" json:"hls_cmaf_url,computed"`
	HlsMpegtsURL        types.String                                                   `tfsdk:"hls_mpegts_url" json:"hls_mpegts_url,computed"`
	IframeURL           types.String                                                   `tfsdk:"iframe_url" json:"iframe_url,computed"`
	Live                types.Bool                                                     `tfsdk:"live" json:"live,computed"`
	PushURL             types.String                                                   `tfsdk:"push_url" json:"push_url,computed"`
	PushURLSrt          types.String                                                   `tfsdk:"push_url_srt" json:"push_url_srt,computed"`
	PushURLWhip         types.String                                                   `tfsdk:"push_url_whip" json:"push_url_whip,computed"`
	RecordingDuration   types.Float64                                                  `tfsdk:"recording_duration" json:"recording_duration,computed"`
	Screenshot          types.String                                                   `tfsdk:"screenshot" json:"screenshot,computed"`
	StartedAtBackup     types.String                                                   `tfsdk:"started_at_backup" json:"started_at_backup,computed"`
	StartedAtPrimary    types.String                                                   `tfsdk:"started_at_primary" json:"started_at_primary,computed"`
	TranscodingSpeed    types.Float64                                                  `tfsdk:"transcoding_speed" json:"transcoding_speed,computed"`
	VideoHeight         types.Float64                                                  `tfsdk:"video_height" json:"video_height,computed"`
	VideoWidth          types.Float64                                                  `tfsdk:"video_width" json:"video_width,computed"`
	TranscodedQualities customfield.List[types.String]                                 `tfsdk:"transcoded_qualities" json:"transcoded_qualities,computed"`
	HTMLOverlays        customfield.NestedObjectList[StreamingStreamHTMLOverlaysModel] `tfsdk:"html_overlays" json:"html_overlays,computed"`
}

func (m StreamingStreamModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m StreamingStreamModel) MarshalJSONForUpdate(state StreamingStreamModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type StreamingStreamStreamModel struct {
	Name                types.String                  `tfsdk:"name" json:"name,required"`
	Active              types.Bool                    `tfsdk:"active" json:"active,computed_optional"`
	AutoRecord          types.Bool                    `tfsdk:"auto_record" json:"auto_record,computed_optional"`
	BroadcastIDs        customfield.List[types.Int64] `tfsdk:"broadcast_ids" json:"broadcast_ids,computed_optional"`
	CdnID               types.Int64                   `tfsdk:"cdn_id" json:"cdn_id,optional"`
	ClientEntityData    types.String                  `tfsdk:"client_entity_data" json:"client_entity_data,optional"`
	ClientUserID        types.Int64                   `tfsdk:"client_user_id" json:"client_user_id,optional"`
	DvrDuration         types.Int64                   `tfsdk:"dvr_duration" json:"dvr_duration,computed_optional"`
	DvrEnabled          types.Bool                    `tfsdk:"dvr_enabled" json:"dvr_enabled,computed_optional"`
	HlsMpegtsEndlistTag types.Bool                    `tfsdk:"hls_mpegts_endlist_tag" json:"hls_mpegts_endlist_tag,computed_optional"`
	HTMLOverlay         types.Bool                    `tfsdk:"html_overlay" json:"html_overlay,computed_optional"`
	LowLatencyEnabled   types.Bool                    `tfsdk:"low_latency_enabled" json:"low_latency_enabled,computed_optional"`
	Projection          types.String                  `tfsdk:"projection" json:"projection,computed_optional"`
	Pull                types.Bool                    `tfsdk:"pull" json:"pull,computed_optional"`
	QualitySetID        types.Int64                   `tfsdk:"quality_set_id" json:"quality_set_id,optional"`
	RecordType          types.String                  `tfsdk:"record_type" json:"record_type,computed_optional"`
	Uri                 types.String                  `tfsdk:"uri" json:"uri,optional"`
}

type StreamingStreamHTMLOverlaysModel struct {
	ID        types.Int64  `tfsdk:"id" json:"id,computed"`
	CreatedAt types.String `tfsdk:"created_at" json:"created_at,computed"`
	StreamID  types.Int64  `tfsdk:"stream_id" json:"stream_id,computed"`
	UpdatedAt types.String `tfsdk:"updated_at" json:"updated_at,computed"`
	URL       types.String `tfsdk:"url" json:"url,computed"`
	Height    types.Int64  `tfsdk:"height" json:"height,computed"`
	Stretch   types.Bool   `tfsdk:"stretch" json:"stretch,computed"`
	Width     types.Int64  `tfsdk:"width" json:"width,computed"`
	X         types.Int64  `tfsdk:"x" json:"x,computed"`
	Y         types.Int64  `tfsdk:"y" json:"y,computed"`
}
