// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_stream

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ resource.ResourceWithConfigValidators = (*StreamingStreamResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:   "Stream ID",
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Description:   "Stream name.\nOften used as a human-readable name for the stream, but can contain any text you wish. The values are not unique and may be repeated.\nExamples:\n- Conference in July\n- Stream #10003\n- Open-Air Camera #31 Backstage\n- 480fd499-2de2-4988-bc1a-a4eebe9818ee",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"cdn_id": schema.Int64Attribute{
				Description:   "ID of custom CDN resource from which the content will be delivered (only if you know what you do)",
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"client_entity_data": schema.StringAttribute{
				Description:   "Custom meta field designed to store your own extra information about a video entity: video source, video id, parameters, etc. We do not use this field in any way when processing the stream. You can store any data in any format (string, json, etc), saved as a text string. Example: ```client_entity_data = '{ \"`seq_id`\": \"1234567890\", \"name\": \"John Doe\", \"iat\": 1516239022 }'```",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"client_user_id": schema.Int64Attribute{
				Description:   "Custom meta field for storing the Identifier in your system. We do not use this field in any way when processing the stream. Example: ```client_user_id = 1001```",
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"quality_set_id": schema.Int64Attribute{
				Description:   "Custom quality set ID for transcoding, if transcoding is required according to your conditions. Look at GET /`quality_sets` method",
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"uri": schema.StringAttribute{
				Description:   "When using PULL method, this is the URL to pull a stream from.\nYou can specify multiple addresses separated by a space (\" \"), so you can organize a backup plan. In this case, the specified addresses will be selected one by one using round robin scheduling. If the first address does not respond, then the next one in the list will be automatically requested, returning to the first and so on in a circle.\nAlso, if the sucessfully working stream stops sending data, then the next one will be selected according to the same scheme.\nAfter 24 hours of inactivity of your streams we will stop PULL-ing, and will switch \"active\" field to \"false\".\nPlease, note that this field is for PULL only, so is not suitable for PUSH. Look at fields \"`push_url`\" and \"`push_url_srt`\" from GET method.",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"active": schema.BoolAttribute{
				Description:   "Stream switch between on and off. This is not an indicator of the status \"stream is receiving and it is LIVE\", but rather an on/off switch.\nWhen stream is switched off, there is no way to process it: PULL is deactivated and PUSH will return an error.\n- true – stream can be processed\n- false – stream is off, and cannot be processed",
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplaceIfConfigured()},
				Default:       booldefault.StaticBool(false),
			},
			"auto_record": schema.BoolAttribute{
				Description:   "Enables autotomatic recording of the stream when it started. So you don't need to call recording manually.\nResult of recording is automatically added to video hosting. For details see the /streams/`start_recording` method and in knowledge base\nValues:\n- true – auto recording is enabled\n- false – auto recording is disabled",
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplaceIfConfigured()},
				Default:       booldefault.StaticBool(false),
			},
			"dvr_duration": schema.Int64Attribute{
				Description:   "DVR duration in seconds if DVR feature is enabled for the stream. So this is duration of how far the user can rewind the live stream.\n`dvr_duration` range is [30...14400].\nMaximum value is 4 hours = 14400 seconds. If you need more, ask the Support Team please.",
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplaceIfConfigured()},
				Default:       int64default.StaticInt64(3600),
			},
			"dvr_enabled": schema.BoolAttribute{
				Description:   "Enables DVR for the stream:\n- true – DVR is enabled\n- false – DVR is disabled",
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplaceIfConfigured()},
				Default:       booldefault.StaticBool(false),
			},
			"hls_mpegts_endlist_tag": schema.BoolAttribute{
				Description:   "Add ```#EXT-X-ENDLIST``` tag within .m3u8 playlist after the last segment of a live stream when broadcast is ended.",
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplaceIfConfigured()},
				Default:       booldefault.StaticBool(true),
			},
			"html_overlay": schema.BoolAttribute{
				Description:   "Switch on mode to insert and display real-time HTML overlay widgets on top of live streams",
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplaceIfConfigured()},
				Default:       booldefault.StaticBool(false),
			},
			"low_latency_enabled": schema.BoolAttribute{
				Description:   "Deprecated, always returns \"true\".\nThe only exception is that the attribute can only be used by clients that have previously used the old stream format.\nThis method is outdated since we've made it easier to manage streams.\nFor your convenience, you no longer need to set this parameter at the stage of creating a stream. Now all streams are prepared in 2 formats simultaniously: Low Latency and Legacy. You can get the desired output format in the attributes \"`dash_url`\", \"`hls_cmaf_url`\", \"`hls_mpegts_url`\". Or use them all at once.\n---\nNote: Links /streams/{id}/playlist.m3u8 are depricated too. Use value of the \"`hls_mpegts_url`\" attribute instead.",
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplaceIfConfigured()},
				Default:       booldefault.StaticBool(true),
			},
			"projection": schema.StringAttribute{
				Description: "Visualization mode for 360° streams, how the stream is rendered in our web player ONLY. If you would like to show video 360° in an external video player, then use parameters of that video player.\nModes:\n- regular – regular “flat” stream\n- vr360 – display stream in 360° mode\n- vr180 – display stream in 180° mode\n- vr360tb – display stream in 3D 360° mode Top-Bottom\nAvailable values: \"regular\", \"vr360\", \"vr180\", \"vr360tb\".",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"regular",
						"vr360",
						"vr180",
						"vr360tb",
					),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIfConfigured()},
				Default:       stringdefault.StaticString("regular"),
			},
			"pull": schema.BoolAttribute{
				Description:   "Indicates if stream is pulled from external server or not. Has two possible\nvalues:\n- true – stream is received by PULL method. Use this when need to get stream from external server by srt, rtmp\\s, hls, dash, etc protocols.\n- false – stream is received by PUSH method. Use this when need to send stream from end-device to our Streaming Platform, i.e. from mobile app or OBS Studio.",
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplaceIfConfigured()},
				Default:       booldefault.StaticBool(false),
			},
			"record_type": schema.StringAttribute{
				Description: "Method of recording a stream. Specifies the source from which the stream will be recorded: original or transcoded.\nTypes:\n- \"origin\" – To record RMTP/SRT/etc original clean media source.\n- \"transcoded\" – To record the output transcoded version of the stream, including overlays, texts, logos, etc. additional media layers.\nAvailable values: \"origin\", \"transcoded\".",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("origin", "transcoded"),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIfConfigured()},
				Default:       stringdefault.StaticString("origin"),
			},
			"broadcast_ids": schema.ListAttribute{
				Description:   "IDs of broadcasts which will include this stream",
				Computed:      true,
				Optional:      true,
				CustomType:    customfield.NewListType[types.Int64](ctx),
				ElementType:   types.Int64Type,
				PlanModifiers: []planmodifier.List{listplanmodifier.RequiresReplaceIfConfigured()},
			},
			"stream": schema.SingleNestedAttribute{
				Computed:   true,
				Optional:   true,
				CustomType: customfield.NewNestedObjectType[StreamingStreamStreamModel](ctx),
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description: "Stream name.\nOften used as a human-readable name for the stream, but can contain any text you wish. The values are not unique and may be repeated.\nExamples:\n- Conference in July\n- Stream #10003\n- Open-Air Camera #31 Backstage\n- 480fd499-2de2-4988-bc1a-a4eebe9818ee",
						Required:    true,
					},
					"active": schema.BoolAttribute{
						Description: "Stream switch between on and off. This is not an indicator of the status \"stream is receiving and it is LIVE\", but rather an on/off switch.\nWhen stream is switched off, there is no way to process it: PULL is deactivated and PUSH will return an error.\n- true – stream can be processed\n- false – stream is off, and cannot be processed",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
					"auto_record": schema.BoolAttribute{
						Description: "Enables autotomatic recording of the stream when it started. So you don't need to call recording manually.\nResult of recording is automatically added to video hosting. For details see the /streams/`start_recording` method and in knowledge base\nValues:\n- true – auto recording is enabled\n- false – auto recording is disabled",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
					"broadcast_ids": schema.ListAttribute{
						Description: "IDs of broadcasts which will include this stream",
						Computed:    true,
						Optional:    true,
						CustomType:  customfield.NewListType[types.Int64](ctx),
						ElementType: types.Int64Type,
					},
					"cdn_id": schema.Int64Attribute{
						Description: "ID of custom CDN resource from which the content will be delivered (only if you know what you do)",
						Optional:    true,
					},
					"client_entity_data": schema.StringAttribute{
						Description: "Custom meta field designed to store your own extra information about a video entity: video source, video id, parameters, etc. We do not use this field in any way when processing the stream. You can store any data in any format (string, json, etc), saved as a text string. Example: ```client_entity_data = '{ \"`seq_id`\": \"1234567890\", \"name\": \"John Doe\", \"iat\": 1516239022 }'```",
						Optional:    true,
					},
					"client_user_id": schema.Int64Attribute{
						Description: "Custom meta field for storing the Identifier in your system. We do not use this field in any way when processing the stream. Example: ```client_user_id = 1001```",
						Optional:    true,
					},
					"dvr_duration": schema.Int64Attribute{
						Description: "DVR duration in seconds if DVR feature is enabled for the stream. So this is duration of how far the user can rewind the live stream.\n`dvr_duration` range is [30...14400].\nMaximum value is 4 hours = 14400 seconds. If you need more, ask the Support Team please.",
						Computed:    true,
						Optional:    true,
						Default:     int64default.StaticInt64(3600),
					},
					"dvr_enabled": schema.BoolAttribute{
						Description: "Enables DVR for the stream:\n- true – DVR is enabled\n- false – DVR is disabled",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
					"hls_mpegts_endlist_tag": schema.BoolAttribute{
						Description: "Add ```#EXT-X-ENDLIST``` tag within .m3u8 playlist after the last segment of a live stream when broadcast is ended.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(true),
					},
					"html_overlay": schema.BoolAttribute{
						Description: "Switch on mode to insert and display real-time HTML overlay widgets on top of live streams",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
					"low_latency_enabled": schema.BoolAttribute{
						Description: "Deprecated, always returns \"true\".\nThe only exception is that the attribute can only be used by clients that have previously used the old stream format.\nThis method is outdated since we've made it easier to manage streams.\nFor your convenience, you no longer need to set this parameter at the stage of creating a stream. Now all streams are prepared in 2 formats simultaniously: Low Latency and Legacy. You can get the desired output format in the attributes \"`dash_url`\", \"`hls_cmaf_url`\", \"`hls_mpegts_url`\". Or use them all at once.\n---\nNote: Links /streams/{id}/playlist.m3u8 are depricated too. Use value of the \"`hls_mpegts_url`\" attribute instead.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(true),
					},
					"projection": schema.StringAttribute{
						Description: "Visualization mode for 360° streams, how the stream is rendered in our web player ONLY. If you would like to show video 360° in an external video player, then use parameters of that video player.\nModes:\n- regular – regular “flat” stream\n- vr360 – display stream in 360° mode\n- vr180 – display stream in 180° mode\n- vr360tb – display stream in 3D 360° mode Top-Bottom\nAvailable values: \"regular\", \"vr360\", \"vr180\", \"vr360tb\".",
						Computed:    true,
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive(
								"regular",
								"vr360",
								"vr180",
								"vr360tb",
							),
						},
						Default: stringdefault.StaticString("regular"),
					},
					"pull": schema.BoolAttribute{
						Description: "Indicates if stream is pulled from external server or not. Has two possible\nvalues:\n- true – stream is received by PULL method. Use this when need to get stream from external server by srt, rtmp\\s, hls, dash, etc protocols.\n- false – stream is received by PUSH method. Use this when need to send stream from end-device to our Streaming Platform, i.e. from mobile app or OBS Studio.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
					"quality_set_id": schema.Int64Attribute{
						Description: "Custom quality set ID for transcoding, if transcoding is required according to your conditions. Look at GET /`quality_sets` method",
						Optional:    true,
					},
					"record_type": schema.StringAttribute{
						Description: "Method of recording a stream. Specifies the source from which the stream will be recorded: original or transcoded.\nTypes:\n- \"origin\" – To record RMTP/SRT/etc original clean media source.\n- \"transcoded\" – To record the output transcoded version of the stream, including overlays, texts, logos, etc. additional media layers.\nAvailable values: \"origin\", \"transcoded\".",
						Computed:    true,
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive("origin", "transcoded"),
						},
						Default: stringdefault.StaticString("origin"),
					},
					"uri": schema.StringAttribute{
						Description: "When using PULL method, this is the URL to pull a stream from.\nYou can specify multiple addresses separated by a space (\" \"), so you can organize a backup plan. In this case, the specified addresses will be selected one by one using round robin scheduling. If the first address does not respond, then the next one in the list will be automatically requested, returning to the first and so on in a circle.\nAlso, if the sucessfully working stream stops sending data, then the next one will be selected according to the same scheme.\nAfter 24 hours of inactivity of your streams we will stop PULL-ing, and will switch \"active\" field to \"false\".\nPlease, note that this field is for PULL only, so is not suitable for PUSH. Look at fields \"`push_url`\" and \"`push_url_srt`\" from GET method.",
						Optional:    true,
					},
				},
			},
			"backup_live": schema.BoolAttribute{
				Description: "State of receiving and transcoding master stream from source by backup server if you pushing stream to \"`backup_push_url`\" or \"`backup_push_url_srt`\".\nDisplays the backup server status of PUSH method only. For PULL a \"live\" field is always used, even when origin servers are switched using round robin scheduling (look \"uri\" field for details).",
				Computed:    true,
			},
			"backup_push_url": schema.StringAttribute{
				Description: "URL to PUSH master stream to our backup server using RTMP/S protocols. Servers for the main and backup streams are distributed geographically.\nMainly sending one stream to main server is enough. But if you need a backup stream, then this is the field to PUSH it.\nTo use RTMPS just manually change the protocol name from \"rtmp://\" to \"rtmps://\".\nThe backup logs are as follows: In PUSH mode, you initiate sending a stream from your machine. If your stream stops or breaks for some reason and it stops coming to the main server, then after 3-10 seconds of waiting the stream will turn off or the backup one will be automatically turned on, if you are pushing it too.",
				Computed:    true,
			},
			"backup_push_url_srt": schema.StringAttribute{
				Description: "URL to PUSH master stream to our backup server using SRT protocol with the same logic of backup-streams",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Datetime of creation in ISO 8601",
				Computed:    true,
			},
			"dash_url": schema.StringAttribute{
				Description: "MPEG-DASH output. URL for transcoded result stream in MPEG-DASH format, with .mpd link.\nLow Latency support: YES.\nThis is CMAF-based MPEG-DASH stream. Encoder and packager dynamically assemble the video stream with fMP4 fragments. Chunks have ±2-4 seconds duration depending on the settings. All chunks for DASH are transferred through CDN using chunk transfer technology, which allows to use all the advantages of low latency delivery of DASH.\n- by default low latency is ±4 sec, because it's stable for almost all last-mile use cases.\n- and its possible to enable ±2 sec for DASH, just ask our Support Team.\n  \nRead more information in the article \"How Low Latency streaming works\" in the Knowledge Base.",
				Computed:    true,
			},
			"finished_at_primary": schema.StringAttribute{
				Description: "Time when the stream ended for the last time. Datetime in ISO 8601.\nAfter restarting the stream, this value is not reset to \"null\", and the time of the last/previous end is always displayed here. That is, when the start time is greater than the end time, it means the current session is still ongoing and the stream has not ended yet.\nIf you want to see all information about acitivity of the stream, you can get it from another method /streaming/statistics/ffprobe. This method shows aggregated activity parameters during a time, when stream was alive and transcoded. Also you can create graphs to see the activity. For example /streaming/statistics/ffprobe?interval=6000&`date_from`=2023-10-01&`date_to`=2023-10-11&`stream_id`=12345",
				Computed:    true,
			},
			"frame_rate": schema.Float64Attribute{
				Description: "Current FPS of the original stream, if stream is transcoding",
				Computed:    true,
			},
			"hls_cmaf_url": schema.StringAttribute{
				Description: "HLS output. URL for transcoded result of stream in HLS CMAF format, with .m3u8 link.\nRecommended for use for all HLS streams.\nLow Latency support: YES.\nThis is CMAF-based HLS stream. Encoder and packager dynamically assemble the video stream with fMP4 fragments. Chunks have ±2-4 seconds duration depending on the settings. All chunks for LL-HLS are transferred through CDN via dividing into parts (small segments ```#EXT-X-PART``` of 0.5-1.0 sec duration), which allows to use all the advantages of low latency delivery of LL-HLS.\n- by default low latency is ±5 sec, because it's stable for almost all last-mile use cases.\n- and its possible to enable ±3 sec for LL-HLS, just ask our Support Team.\n  \nIt is also possible to use modifier-attributes, which are described in the \"`hls_mpegts_url`\" field above.\nIf you need to get MPEGTS (.ts) chunks, look at the attribute \"`hls_mpegts_url`\".\n  \nRead more information in the article \"How Low Latency streaming works\" in the Knowledge Base.",
				Computed:    true,
			},
			"hls_mpegts_url": schema.StringAttribute{
				Description: "HLS output for legacy devices. URL for transcoded result of stream in HLS MPEGTS (.ts) format, with .m3u8 link.\nLow Latency support: NO.\nSome legacy devices or software may require MPEGTS (.ts) segments as a format for streaming, so we provide this options keeping backward compatibility with any of your existing workflows. For other cases it's better to use \"`hls_cmaf_url`\" instead.\nYou can use this legacy HLSv6 format based on MPEGTS segmenter in parallel with main HLS CMAF. Both formats are sharing same segments size, manifest length (DVR), etc.\n  \nIt is also possible to use additional modifier-attributes:\n- ?`get_duration_sec`=true – Adds the real segment duration in seconds to chunk requests. A chunk duration will be automatically added to a chunk request string with the \"`duration_sec`\" attribute. The value is an integer for a length multiple of whole seconds, or a fractional number separated by a dot for chunks that are not multiples of seconds. This attribute allows you to determine duration in seconds at the level of analyzing the logs of CDN requests and compare it with file size (so to use it in your analytics).\nSuch modifier attributes are applied manually and added to the link obtained from this field. I.e. ```?`get_duration_sec`=true```\nExample:\n`https://demo.gvideo.io/mpegts/`2675_19146`/`master_mpegts`.m3u8?`get_duration_sec`=true`\n```\n#EXTM3U\n#EXT-X-VERSION:6\n#EXT-X-TARGETDURATION:2\n...\n#EXTINF:2.000000,\n#EXT-X-PROGRAM-DATE-TIME:2025-08-14T08:15:00\nseg1.ts?`duration_sec`=2\n...\n```",
				Computed:    true,
			},
			"iframe_url": schema.StringAttribute{
				Description: "A URL to a built-in HTML web player with the stream inside. It can be inserted into an iframe on your website and the video will automatically play in all browsers.\nPlease, remember that transcoded streams from \"`hls_cmaf_url`\" with .m3u8 at the end, and from \"`dash_url`\" with .mpd at the end are to be played inside video players only. For example: AVplayer on iOS, Exoplayer on Android, HTML web player in browser, etc. General bowsers like Chrome, Firefox, etc cannot play transcoded streams with .m3u8 and .mpd at the end. The only exception is Safari, which can only play Apple's HLS .m3u8 format with limits.\nThat's why you may need to use this HTML web player. Please, look Knowledge Base for details.\nExample of usage on a web page:\n<iframe width=\"560\" height=\"315\" src=\"https://player.gvideo.co/streams/`2675_201693`\" frameborder=\"0\" allow=\"autoplay; encrypted-media\" allowfullscreen></iframe>",
				Computed:    true,
			},
			"live": schema.BoolAttribute{
				Description: "State of receiving and transcoding master stream from source by main server",
				Computed:    true,
			},
			"push_url": schema.StringAttribute{
				Description: "URL to PUSH master stream to our main server using RTMP and RTMPS protocols.\nTo use RTMPS just manually change the protocol name from \"rtmp://\" to \"rtmps://\".\nIf you see an error like \"invalid SSL certificate\" try the following:\n- Make sure the push URL is correct, and it contains \"rtmps://\".\n- If the URL looks correct but you still get an SSL error, try specifying the port 443 in the URL. Here’s an example: rtmps://vp-push.domain.com:443/in/stream?key.\n- If you're still having trouble, then your encoder may not support RTMPS. Double-check the documentation for your encoder.\nFor advanced customers only: For your complexly distributed broadcast systems, it is also possible to additionally output an array of multi-regional ingestion points for manual selection from them. To activate this mode, contact your manager or the Support Team to activate the \"`multi_region_push_urls`\" attibute.\nBut if you clearly don’t understand why you need this, then it’s best to use the default single URL in the \"`push_url`\" attribute.",
				Computed:    true,
			},
			"push_url_srt": schema.StringAttribute{
				Description: "URL to PUSH master stream to our main server using SRT protocol.\nUse only 1 protocol of sending a master stream: either only SRT (`push_url_srt`), or only RTMP (`push_url`).",
				Computed:    true,
			},
			"push_url_whip": schema.StringAttribute{
				Description: "URL to PUSH WebRTC stream to our server using WHIP protocol.\n  \n**WebRTC WHIP to LL-HLS and DASH**\nVideo Streaming supports WebRTC HTTP Ingest Protocol (WHIP), and WebRTC to HLS/DASH converter. As a result you can stream from web broswers natively.\n**WebRTC WHIP server**\nWe have dedicated WebRTC WHIP servers in our infrastructure. WebRTC WHIP server organizes both signaling and receives video data.\nSignaling is a term to describe communication between WebRTC endpoints, needed to initiate and maintain a session. WHIP is an open specification for a simple signaling protocol for starting WebRTC sessions in an outgoing direction, (i.e., streaming from your device).\n**WebRTC stream encoding parameters**\nAt least one video and audio track both must be present in the stream:\n- Video must be encoded with H.264.\n- Audio must be encoded with OPUS.\nNote. Specifically for WebRTC mode a method of constant transcoding with an initial given resolution is used. This means that if WebRTC in the end-user's browser decides to reduce the quality or resolution of the master stream (to let say 360p) due to restrictions on the end-user's device (network conditions, CPU consumption, etc.), the transcoder will still continue to transcode the reduced stream to the initial resolution (let say 1080p ABR). When the restrictions on the end-user's device are removed, quiality will improve again.\n**WebRTC WHIP Client**\nWe provide a convenient WebRTC WHIP library for working in browsers. You can use our library, or any other you prefer.\nSimple example of usage is here: https://stackblitz.com/edit/stackblitz-starters-j2r9ar?file=index.html\nAlso try to use the feature in UI of the Customer Portal. In the Streaming section inside the settings of a specific live stream, a new section \"Quick start in browser\" has been added.\nMore information in the Product Documentation on the website.",
				Computed:    true,
			},
			"recording_duration": schema.Float64Attribute{
				Description: "Duration of current recording in seconds if recording is enabled for the stream",
				Computed:    true,
			},
			"screenshot": schema.StringAttribute{
				Description: "An instant screenshot taken from a live stream, and available as a static JPEG image. Resolution 1080 pixels wide, or less if the original stream has a lower resolution.\nScreenshot is taken every 10 seconds while the stream is live. This field contains a link to the last screenshot created by the system. Screenshot history is not stored, so if you need a series of screenshots over time, then download them.",
				Computed:    true,
			},
			"started_at_backup": schema.StringAttribute{
				Description: "Time of the last session when backup server started receiving the stream. Datetime in ISO 8601",
				Computed:    true,
			},
			"started_at_primary": schema.StringAttribute{
				Description: "Time of the last session when main server started receiving the stream. Datetime in ISO 8601.\nThis means that if the stream was started 1 time, then here will be the time it was started. If the stream was started several times, or restarted on your side, then only the time of the last session is displayed here.",
				Computed:    true,
			},
			"transcoding_speed": schema.Float64Attribute{
				Description: "Speed of transcoding the stream.\nMainly it must be 1.0 for real-time processing. May be less than 1.0 if your stream has problems in delivery due to your local internet provider's conditions, or the stream does not meet stream inbound requirements. See Knowledge Base for details.",
				Computed:    true,
			},
			"video_height": schema.Float64Attribute{
				Description: "Current height of frame of the original stream, if stream is transcoding",
				Computed:    true,
			},
			"video_width": schema.Float64Attribute{
				Description: "Current width of frame of the original stream, if stream is transcoding",
				Computed:    true,
			},
			"transcoded_qualities": schema.ListAttribute{
				Description: "Array of qualities to which live stream is transcoded",
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
			"html_overlays": schema.ListNestedAttribute{
				Description: "Array of HTML overlay widgets",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[StreamingStreamHTMLOverlaysModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "ID of the overlay",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "Datetime of creation in ISO 8601",
							Computed:    true,
						},
						"stream_id": schema.Int64Attribute{
							Description: "ID of a stream to which it is attached",
							Computed:    true,
						},
						"updated_at": schema.StringAttribute{
							Description: "Datetime of last update in ISO 8601",
							Computed:    true,
						},
						"url": schema.StringAttribute{
							Description: "Valid http/https URL to an HTML page/widget",
							Computed:    true,
						},
						"height": schema.Int64Attribute{
							Description: "Height of the widget",
							Computed:    true,
						},
						"stretch": schema.BoolAttribute{
							Description: `Switch of auto scaling the widget. Must not be used as "true" simultaneously with the coordinate installation method (w, h, x, y).`,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
						"width": schema.Int64Attribute{
							Description: "Width of the widget",
							Computed:    true,
						},
						"x": schema.Int64Attribute{
							Description: "Coordinate of left upper corner",
							Computed:    true,
						},
						"y": schema.Int64Attribute{
							Description: "Coordinate of left upper corner",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (r *StreamingStreamResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *StreamingStreamResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
