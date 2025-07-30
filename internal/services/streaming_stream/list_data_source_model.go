// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_stream

import (
	"context"

	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/gcore-go/streaming"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type StreamingStreamsItemsListDataSourceEnvelope struct {
	Items customfield.NestedObjectList[StreamingStreamsItemsDataSourceModel] `json:"items,computed"`
}

type StreamingStreamsDataSourceModel struct {
	WithBroadcasts types.Int64                                                        `tfsdk:"with_broadcasts" query:"with_broadcasts,optional"`
	MaxItems       types.Int64                                                        `tfsdk:"max_items"`
	Items          customfield.NestedObjectList[StreamingStreamsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *StreamingStreamsDataSourceModel) toListParams(_ context.Context) (params streaming.StreamListParams, diags diag.Diagnostics) {
	params = streaming.StreamListParams{}

	if !m.WithBroadcasts.IsNull() {
		params.WithBroadcasts = param.NewOpt(m.WithBroadcasts.ValueInt64())
	}

	return
}

type StreamingStreamsItemsDataSourceModel struct {
	Name                types.String                                                              `tfsdk:"name" json:"name,computed"`
	ID                  types.Int64                                                               `tfsdk:"id" json:"id,computed"`
	Active              types.Bool                                                                `tfsdk:"active" json:"active,computed"`
	AutoRecord          types.Bool                                                                `tfsdk:"auto_record" json:"auto_record,computed"`
	BackupLive          types.Bool                                                                `tfsdk:"backup_live" json:"backup_live,computed"`
	BackupPushURL       types.String                                                              `tfsdk:"backup_push_url" json:"backup_push_url,computed"`
	BackupPushURLSrt    types.String                                                              `tfsdk:"backup_push_url_srt" json:"backup_push_url_srt,computed"`
	BroadcastIDs        customfield.List[types.Int64]                                             `tfsdk:"broadcast_ids" json:"broadcast_ids,computed"`
	CdnID               types.Int64                                                               `tfsdk:"cdn_id" json:"cdn_id,computed"`
	ClientEntityData    types.String                                                              `tfsdk:"client_entity_data" json:"client_entity_data,computed"`
	ClientUserID        types.Int64                                                               `tfsdk:"client_user_id" json:"client_user_id,computed"`
	CreatedAt           types.String                                                              `tfsdk:"created_at" json:"created_at,computed"`
	DashURL             types.String                                                              `tfsdk:"dash_url" json:"dash_url,computed"`
	DvrDuration         types.Int64                                                               `tfsdk:"dvr_duration" json:"dvr_duration,computed"`
	DvrEnabled          types.Bool                                                                `tfsdk:"dvr_enabled" json:"dvr_enabled,computed"`
	FinishedAtPrimary   types.String                                                              `tfsdk:"finished_at_primary" json:"finished_at_primary,computed"`
	FrameRate           types.Float64                                                             `tfsdk:"frame_rate" json:"frame_rate,computed"`
	HlsCmafURL          types.String                                                              `tfsdk:"hls_cmaf_url" json:"hls_cmaf_url,computed"`
	HlsMpegtsEndlistTag types.Bool                                                                `tfsdk:"hls_mpegts_endlist_tag" json:"hls_mpegts_endlist_tag,computed"`
	HlsMpegtsURL        types.String                                                              `tfsdk:"hls_mpegts_url" json:"hls_mpegts_url,computed"`
	HTMLOverlay         types.Bool                                                                `tfsdk:"html_overlay" json:"html_overlay,computed"`
	HTMLOverlays        customfield.NestedObjectList[StreamingStreamsHTMLOverlaysDataSourceModel] `tfsdk:"html_overlays" json:"html_overlays,computed"`
	IframeURL           types.String                                                              `tfsdk:"iframe_url" json:"iframe_url,computed"`
	Live                types.Bool                                                                `tfsdk:"live" json:"live,computed"`
	LowLatencyEnabled   types.Bool                                                                `tfsdk:"low_latency_enabled" json:"low_latency_enabled,computed"`
	Projection          types.String                                                              `tfsdk:"projection" json:"projection,computed"`
	Pull                types.Bool                                                                `tfsdk:"pull" json:"pull,computed"`
	PushURL             types.String                                                              `tfsdk:"push_url" json:"push_url,computed"`
	PushURLSrt          types.String                                                              `tfsdk:"push_url_srt" json:"push_url_srt,computed"`
	PushURLWhip         types.String                                                              `tfsdk:"push_url_whip" json:"push_url_whip,computed"`
	QualitySetID        types.Int64                                                               `tfsdk:"quality_set_id" json:"quality_set_id,computed"`
	RecordType          types.String                                                              `tfsdk:"record_type" json:"record_type,computed"`
	RecordingDuration   types.Float64                                                             `tfsdk:"recording_duration" json:"recording_duration,computed"`
	Screenshot          types.String                                                              `tfsdk:"screenshot" json:"screenshot,computed"`
	StartedAtBackup     types.String                                                              `tfsdk:"started_at_backup" json:"started_at_backup,computed"`
	StartedAtPrimary    types.String                                                              `tfsdk:"started_at_primary" json:"started_at_primary,computed"`
	TranscodedQualities customfield.List[types.String]                                            `tfsdk:"transcoded_qualities" json:"transcoded_qualities,computed"`
	TranscodingSpeed    types.Float64                                                             `tfsdk:"transcoding_speed" json:"transcoding_speed,computed"`
	Uri                 types.String                                                              `tfsdk:"uri" json:"uri,computed"`
	VideoHeight         types.Float64                                                             `tfsdk:"video_height" json:"video_height,computed"`
	VideoWidth          types.Float64                                                             `tfsdk:"video_width" json:"video_width,computed"`
}

type StreamingStreamsHTMLOverlaysDataSourceModel struct {
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
