// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_video

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ resource.ResourceWithConfigValidators = (*StreamingVideoResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"video_id": schema.Int64Attribute{
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"video": schema.SingleNestedAttribute{
				Computed:   true,
				Optional:   true,
				CustomType: customfield.NewNestedObjectType[StreamingVideoVideoModel](ctx),
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description: "Video name",
						Required:    true,
					},
					"auto_transcribe_audio_language": schema.StringAttribute{
						Description: "Automatic creation of subtitles by transcribing the audio track.\nValues:\n- disable – Do not transcribe.\n- auto – Automatically detects the activation of the option based on the settings in your account. If generation is activated, then automatic language detection while transcribing.\n- \\ – Transcribe from specific language. Can be used to specify the exact language spoken in the audio track, or when auto language detection fails. Language is set by 3-letter language code according to ISO-639-2 (bibliographic code). List of languages is available in ```audio_language``` attribute of API POST /streaming/ai/transcribe .\nExample:\n```\n`auto_transcribe_audio_language`: \"auto\"\n`auto_transcribe_audio_language`: \"ger\"\n```\nMore details:\n- List of AI tasks – API [GET /streaming/ai/tasks](https://api.gcore.com/docs/streaming/docs/api-reference/streaming/ai/get-ai-task-result)\n- Add subtitles to an exist video – API [POST /streaming/videos/{`video_id`}/subtitles](https://api.gcore.com/docs/streaming/docs/api-reference/streaming/subtitles/add-subtitle).\nAvailable values: \"disable\", \"auto\", \"<language_code>\".",
						Computed:    true,
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive(
								"disable",
								"auto",
								"<language_code>",
							),
						},
						Default: stringdefault.StaticString("auto"),
					},
					"auto_translate_subtitles_language": schema.StringAttribute{
						Description: "Automatic translation of auto-transcribed subtitles to the specified language(s). Can be used both together with ```auto_transcribe_audio_language``` option only.\nUse it when you want to make automatic subtitles in languages other than the original language in audio.\nValues:\n- disable – Do not translate.\n- default – There are 3 default languages: eng,fre,ger\n- \\ – Explicit language to translate to, or list of languages separated by a comma. Look at list of available languages in description of AI ASR task creation.\nIf several languages are specified for translation, a separate subtitle will be generated for each language.\nExample:\n```\n`auto_translate_subtitles_language`: default\n`auto_translate_subtitles_language`: eng,fre,ger\n```\n  \nPlease note that subtitle translation is done separately and after transcription. Thus separate AI-tasks are created for translation.\nAvailable values: \"disable\", \"default\", \"<language_codes,>\".",
						Computed:    true,
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive(
								"disable",
								"default",
								"<language_codes,>",
							),
						},
						Default: stringdefault.StaticString("disable"),
					},
					"client_user_id": schema.Int64Attribute{
						Description: "Custom field where you can specify user ID in your system",
						Optional:    true,
					},
					"clip_duration_seconds": schema.Int64Attribute{
						Description: "The length of the trimmed segment to transcode, instead of the entire length of the video. Is only used in conjunction with specifying the start of a segment. Transcoding duration is a number in seconds.",
						Optional:    true,
					},
					"clip_start_seconds": schema.Int64Attribute{
						Description: "If you want to transcode only a trimmed segment of a video instead of entire length if the video, then you can provide timecodes of starting point and duration of a segment to process. Start encoding from is a number in seconds.",
						Optional:    true,
					},
					"custom_iframe_url": schema.StringAttribute{
						Description: "Deprecated.\nCustom URL of IFrame for video player to be used in share panel in\nplayer. Auto generated IFrame URL provided by default",
						Optional:    true,
					},
					"description": schema.StringAttribute{
						Description: "Video details; not visible to the end-users",
						Optional:    true,
					},
					"directory_id": schema.Int64Attribute{
						Description: "ID of the directory where the video should be uploaded. (beta)",
						Optional:    true,
					},
					"origin_http_headers": schema.StringAttribute{
						Description: "Authorization HTTP request header. Will be used as credentials to authenticate a request to download a file (specified in \"`origin_url`\" parameter) on an external server.\nSyntax: ```Authorization:  ```\nExamples:\n- \"`origin_http_headers`\": \"Authorization: Basic ...\"\n- \"`origin_http_headers`\": \"Authorization: Bearer ...\"\n- \"`origin_http_headers`\": \"Authorization: APIKey ...\"\nExample of usage when downloading a file from Google Drive:\n```\nPOST https://api.gcore.com/streaming/videos\n\"video\": {\n\"name\": \"IBC 2024 intro.mp4\",\n\"`origin_url`\": \"https://www.googleapis.com/drive/v3/files/...?alt=media\",\n\"`origin_http_headers`\": \"Authorization: Bearer ABC\"\n}\n```",
						Optional:    true,
					},
					"origin_url": schema.StringAttribute{
						Description: "URL to an original file which you want to copy from external storage.\nIf specified, system will download the file and will use it as video source for transcoding.",
						Optional:    true,
					},
					"poster": schema.StringAttribute{
						Description: "Poster is your own static image which can be displayed before the video starts.\nAfter uploading the video, the system will automatically create several screenshots (they will be stored in \"screenshots\" attribute) from which you can select an default screenshot.\nThis \"poster\" field is for uploading your own image. Also use attribute \"`screenshot_id`\" to select poster as a default screnshot.\nAttribute accepts single image as base64-encoded string [(RFC 2397 – The \"data\" URL scheme)](https://www.rfc-editor.org/rfc/rfc2397). In format: ```data:[];base64,```\nMIME-types are image/jpeg, image/webp, and image/png and file sizes up to 1Mb.\nExamples:\n- ```data:image/jpeg;base64,/9j/4AA...qf/2Q==```\n- ```data:image/png;base64,iVBORw0KGg...ggg==```\n- ```data:image/webp;base64,UklGRt.../DgAAAAA```",
						Optional:    true,
					},
					"priority": schema.Int64Attribute{
						Description: "Priority allows you to adjust the urgency of processing some videos before others in your account, if your algorithm requires it. For example, when there are very urgent video and some regular ones that can wait in the queue.\nValue range, integer [-10..10]. -10 is the lowest down-priority, 10 is the highest up-priority. Default priority is 0.",
						Computed:    true,
						Optional:    true,
						Default:     int64default.StaticInt64(0),
					},
					"projection": schema.StringAttribute{
						Description: "Deprecated.\nRegulates the video format:\n\n* **regular**\n  — plays the video as usual\n* **vr360**\n  — plays the video in 360 degree mode\n* **vr180** — plays the video in 180 degree\n  mode\n* **vr360tb** — plays the video in 3D\n  360 degree mode Top-Bottom.\n\n  \n Default is regular",
						Optional:    true,
					},
					"quality_set_id": schema.Int64Attribute{
						Description: "Custom quality set ID for transcoding, if transcoding is required according to your conditions. Look at GET /`quality_sets` method",
						Optional:    true,
					},
					"remote_poster_url": schema.StringAttribute{
						Description: "Poster URL to download from external resource, instead of uploading via \"poster\" attribute.\nIt has the same restrictions as \"poster\" attribute.",
						Optional:    true,
					},
					"remove_poster": schema.BoolAttribute{
						Description: "Set it to true to remove poster",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
					"screenshot_id": schema.Int64Attribute{
						Description: "Default screenshot index.\nSpecify an ID from the \"screenshots\" array, so that the URL of the required screenshot appears in the \"screenshot\" attribute as the default screenshot. By default 5 static screenshots will be taken from different places in the video after transcoding. If the video is short, there may be fewer screenshots.\nCounting from 0. A value of -1 sets the default screenshot to the URL of your own image from the \"poster\" attribute.\nLook at \"screenshot\" attribute in GET /videos/{`video_id`} for details.",
						Computed:    true,
						Optional:    true,
						Default:     int64default.StaticInt64(0),
					},
					"share_url": schema.StringAttribute{
						Description: "Deprecated.\nCustom URL or iframe displayed in the link field when a user clicks\non a sharing button in player. If empty, the link field and social\nnetwork sharing is disabled",
						Optional:    true,
					},
					"source_bitrate_limit": schema.BoolAttribute{
						Description: "The option allows you to set the video transcoding rule so that the output bitrate in ABR ladder is not exceeding the bitrate of the original video.\n  \nThis option is for advanced users only.\n  \nBy default ```source_bitrate_limit: true``` this option allows you to have the output bitrate not more than in the original video, thus to transcode video faster and to deliver it to end-viewers faster as well. At the same time, the quality will be similar to the original.\nIf for some reason you need more byte-space in the output quality when encoding, you can set this option to ```source_bitrate_limit: false```. Then, when transcoding, the quality ceiling will be raised from the bitrate of the original video to the maximum possible limit specified in our the Product Documentation.\nFor example, this may be needed when:\n- to improve the visual quality parameters using PSNR, SSIM, VMAF metrics,\n- to improve the picture quality on dynamic scenes,\n- etc.\nThe option is applied only at the video creation stage and cannot be changed later. If you want to re-transcode the video using new value, then you need to create and upload a new video only.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(true),
					},
				},
				PlanModifiers: []planmodifier.Object{objectplanmodifier.RequiresReplaceIfConfigured()},
			},
			"client_user_id": schema.Int64Attribute{
				Description: "Custom field where you can specify user ID in your system",
				Optional:    true,
			},
			"clip_duration_seconds": schema.Int64Attribute{
				Description: "The length of the trimmed segment to transcode, instead of the entire length of the video. Is only used in conjunction with specifying the start of a segment. Transcoding duration is a number in seconds.",
				Optional:    true,
			},
			"clip_start_seconds": schema.Int64Attribute{
				Description: "If you want to transcode only a trimmed segment of a video instead of entire length if the video, then you can provide timecodes of starting point and duration of a segment to process. Start encoding from is a number in seconds.",
				Optional:    true,
			},
			"custom_iframe_url": schema.StringAttribute{
				Description: "Deprecated.\nCustom URL of IFrame for video player to be used in share panel in\nplayer. Auto generated IFrame URL provided by default",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "Video details; not visible to the end-users",
				Optional:    true,
			},
			"directory_id": schema.Int64Attribute{
				Description: "ID of the directory where the video should be uploaded. (beta)",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "Video name",
				Optional:    true,
			},
			"origin_http_headers": schema.StringAttribute{
				Description: "Authorization HTTP request header. Will be used as credentials to authenticate a request to download a file (specified in \"`origin_url`\" parameter) on an external server.\nSyntax: ```Authorization:  ```\nExamples:\n- \"`origin_http_headers`\": \"Authorization: Basic ...\"\n- \"`origin_http_headers`\": \"Authorization: Bearer ...\"\n- \"`origin_http_headers`\": \"Authorization: APIKey ...\"\nExample of usage when downloading a file from Google Drive:\n```\nPOST https://api.gcore.com/streaming/videos\n\"video\": {\n\"name\": \"IBC 2024 intro.mp4\",\n\"`origin_url`\": \"https://www.googleapis.com/drive/v3/files/...?alt=media\",\n\"`origin_http_headers`\": \"Authorization: Bearer ABC\"\n}\n```",
				Optional:    true,
			},
			"origin_url": schema.StringAttribute{
				Description: "URL to an original file which you want to copy from external storage.\nIf specified, system will download the file and will use it as video source for transcoding.",
				Optional:    true,
			},
			"poster": schema.StringAttribute{
				Description: "Poster is your own static image which can be displayed before the video starts.\nAfter uploading the video, the system will automatically create several screenshots (they will be stored in \"screenshots\" attribute) from which you can select an default screenshot.\nThis \"poster\" field is for uploading your own image. Also use attribute \"`screenshot_id`\" to select poster as a default screnshot.\nAttribute accepts single image as base64-encoded string [(RFC 2397 – The \"data\" URL scheme)](https://www.rfc-editor.org/rfc/rfc2397). In format: ```data:[];base64,```\nMIME-types are image/jpeg, image/webp, and image/png and file sizes up to 1Mb.\nExamples:\n- ```data:image/jpeg;base64,/9j/4AA...qf/2Q==```\n- ```data:image/png;base64,iVBORw0KGg...ggg==```\n- ```data:image/webp;base64,UklGRt.../DgAAAAA```",
				Optional:    true,
			},
			"projection": schema.StringAttribute{
				Description: "Deprecated.\nRegulates the video format:\n\n* **regular**\n  — plays the video as usual\n* **vr360**\n  — plays the video in 360 degree mode\n* **vr180** — plays the video in 180 degree\n  mode\n* **vr360tb** — plays the video in 3D\n  360 degree mode Top-Bottom.\n\n  \n Default is regular",
				Optional:    true,
			},
			"quality_set_id": schema.Int64Attribute{
				Description: "Custom quality set ID for transcoding, if transcoding is required according to your conditions. Look at GET /`quality_sets` method",
				Optional:    true,
			},
			"remote_poster_url": schema.StringAttribute{
				Description: "Poster URL to download from external resource, instead of uploading via \"poster\" attribute.\nIt has the same restrictions as \"poster\" attribute.",
				Optional:    true,
			},
			"share_url": schema.StringAttribute{
				Description: "Deprecated.\nCustom URL or iframe displayed in the link field when a user clicks\non a sharing button in player. If empty, the link field and social\nnetwork sharing is disabled",
				Optional:    true,
			},
			"auto_transcribe_audio_language": schema.StringAttribute{
				Description: "Automatic creation of subtitles by transcribing the audio track.\nValues:\n- disable – Do not transcribe.\n- auto – Automatically detects the activation of the option based on the settings in your account. If generation is activated, then automatic language detection while transcribing.\n- \\ – Transcribe from specific language. Can be used to specify the exact language spoken in the audio track, or when auto language detection fails. Language is set by 3-letter language code according to ISO-639-2 (bibliographic code). List of languages is available in ```audio_language``` attribute of API POST /streaming/ai/transcribe .\nExample:\n```\n`auto_transcribe_audio_language`: \"auto\"\n`auto_transcribe_audio_language`: \"ger\"\n```\nMore details:\n- List of AI tasks – API [GET /streaming/ai/tasks](https://api.gcore.com/docs/streaming/docs/api-reference/streaming/ai/get-ai-task-result)\n- Add subtitles to an exist video – API [POST /streaming/videos/{`video_id`}/subtitles](https://api.gcore.com/docs/streaming/docs/api-reference/streaming/subtitles/add-subtitle).\nAvailable values: \"disable\", \"auto\", \"<language_code>\".",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"disable",
						"auto",
						"<language_code>",
					),
				},
				Default: stringdefault.StaticString("auto"),
			},
			"auto_translate_subtitles_language": schema.StringAttribute{
				Description: "Automatic translation of auto-transcribed subtitles to the specified language(s). Can be used both together with ```auto_transcribe_audio_language``` option only.\nUse it when you want to make automatic subtitles in languages other than the original language in audio.\nValues:\n- disable – Do not translate.\n- default – There are 3 default languages: eng,fre,ger\n- \\ – Explicit language to translate to, or list of languages separated by a comma. Look at list of available languages in description of AI ASR task creation.\nIf several languages are specified for translation, a separate subtitle will be generated for each language.\nExample:\n```\n`auto_translate_subtitles_language`: default\n`auto_translate_subtitles_language`: eng,fre,ger\n```\n  \nPlease note that subtitle translation is done separately and after transcription. Thus separate AI-tasks are created for translation.\nAvailable values: \"disable\", \"default\", \"<language_codes,>\".",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"disable",
						"default",
						"<language_codes,>",
					),
				},
				Default: stringdefault.StaticString("disable"),
			},
			"priority": schema.Int64Attribute{
				Description: "Priority allows you to adjust the urgency of processing some videos before others in your account, if your algorithm requires it. For example, when there are very urgent video and some regular ones that can wait in the queue.\nValue range, integer [-10..10]. -10 is the lowest down-priority, 10 is the highest up-priority. Default priority is 0.",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
			},
			"remove_poster": schema.BoolAttribute{
				Description: "Set it to true to remove poster",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"screenshot_id": schema.Int64Attribute{
				Description: "Default screenshot index.\nSpecify an ID from the \"screenshots\" array, so that the URL of the required screenshot appears in the \"screenshot\" attribute as the default screenshot. By default 5 static screenshots will be taken from different places in the video after transcoding. If the video is short, there may be fewer screenshots.\nCounting from 0. A value of -1 sets the default screenshot to the URL of your own image from the \"poster\" attribute.\nLook at \"screenshot\" attribute in GET /videos/{`video_id`} for details.",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
			},
			"source_bitrate_limit": schema.BoolAttribute{
				Description: "The option allows you to set the video transcoding rule so that the output bitrate in ABR ladder is not exceeding the bitrate of the original video.\n  \nThis option is for advanced users only.\n  \nBy default ```source_bitrate_limit: true``` this option allows you to have the output bitrate not more than in the original video, thus to transcode video faster and to deliver it to end-viewers faster as well. At the same time, the quality will be similar to the original.\nIf for some reason you need more byte-space in the output quality when encoding, you can set this option to ```source_bitrate_limit: false```. Then, when transcoding, the quality ceiling will be raised from the bitrate of the original video to the maximum possible limit specified in our the Product Documentation.\nFor example, this may be needed when:\n- to improve the visual quality parameters using PSNR, SSIM, VMAF metrics,\n- to improve the picture quality on dynamic scenes,\n- etc.\nThe option is applied only at the video creation stage and cannot be changed later. If you want to re-transcode the video using new value, then you need to create and upload a new video only.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(true),
			},
			"ad_id": schema.Int64Attribute{
				Description: "ID of ad that should be shown. If empty the default ad is show. If there is no default ad, no ad is shownю",
				Computed:    true,
			},
			"cdn_views": schema.Int64Attribute{
				Description: "Total number of video views. It is calculated based on the analysis of all views, no matter in which player.",
				Computed:    true,
			},
			"client_id": schema.Int64Attribute{
				Description: "Client ID",
				Computed:    true,
			},
			"dash_url": schema.StringAttribute{
				Description: "A URL to a master playlist MPEG-DASH (master.mpd) with CMAF or WebM based chunks.\nChunk type will be selected automatically for each quality:\n- CMAF for H264 and H265 codecs.\n- WebM for AV1 codec.\n  \nThis URL is a link to the main manifest. But you can also manually specify suffix-options that will allow you to change the manifest to your request:\n```/videos/{`client_id`}_{slug}/master[-min-N][-max-N][-(h264|hevc|av1)].mpd```\nList of suffix-options:\n- [-min-N] – ABR soft limitation of qualities from below.\n- [-max-N] – ABR soft limitation of qualities from above.\n- [-(h264|hevc|av1) – Video codec soft limitation. Applicable if the video was transcoded into multiple codecs H264, H265 and AV1 at once, but you want to return just 1 video codec in a manifest. Read the Product Documentation for details.\nRead more what is ABR soft-limiting in the \"`hls_url`\" field above.\n  \nCaution. Solely master.mpd is officially documented and intended for your use. Any additional internal manifests, sub-manifests, parameters, chunk names, file extensions, and related components are internal infrastructure entities. These may undergo modifications without prior notice, in any manner or form. It is strongly advised not to store them in your database or cache them on your end.",
				Computed:    true,
			},
			"duration": schema.Int64Attribute{
				Description: "Video duration in milliseconds. May differ from \"`origin_video_duration`\" value if the video was uploaded with clipping through the parameters \"`clip_start_seconds`\" and \"`clip_duration_seconds`\"",
				Computed:    true,
			},
			"error": schema.StringAttribute{
				Description: `Video processing error text will be saved here if "status: error"`,
				Computed:    true,
			},
			"hls_cmaf_url": schema.StringAttribute{
				Description: "A URL to a master playlist HLS (master-cmaf.m3u8) with CMAF-based chunks. Chunks are in fMP4 container. It's a code-agnostic container, which allows to use any like H264, H265, AV1, etc.\n  \nIt is possible to use the same suffix-options as described in the \"`hls_url`\" attribute.\n  \nCaution. Solely master.m3u8 (and master[-options].m3u8) is officially documented and intended for your use. Any additional internal manifests, sub-manifests, parameters, chunk names, file extensions, and related components are internal infrastructure entities. These may undergo modifications without prior notice, in any manner or form. It is strongly advised not to store them in your database or cache them on your end.",
				Computed:    true,
			},
			"hls_url": schema.StringAttribute{
				Description: "A URL to a master playlist HLS (master.m3u8).\nChunk type will be selected automatically:\n- TS if your video was encoded to H264 only.\n- CMAF if your video was encoded additionally to H265 and/or AV1 codecs (as Apple does not support these codecs over MPEG TS, and they are not standardized in TS-container).\n  \nYou can also manually specify suffix-options that will allow you to change the manifest to your request:\n```/videos/{`client_id`}_{`video_slug`}/master[-cmaf][-min-N][-max-N][-img][-(h264|hevc|av1)].m3u8```\nList of suffix-options:\n- [-cmaf] – getting HLS CMAF version of the manifest. Look at the ```hls_cmaf_url``` field.\n- [-min-N] – ABR soft limitation of qualities from below.\n- [-max-N] – ABR soft limitation of qualities from above.\n- [-img] – Roku trick play: to add tiles directly into .m3u8 manifest. Read the Product Documentation for details.\n- [-(h264|hevc|av1) – Video codec soft limitation. Applicable if the video was transcoded into multiple codecs H264, H265 and AV1 at once, but you want to return just 1 video codec in a manifest. Read the Product Documentation for details.\nABR soft-limiting: Soft limitation of the list of qualities allows you to return not the entire list of transcoded qualities for a video, but only those you need. For more details look at the Product Documentation.\nFor example, the video is available in 7 qualities from 360p to 4K, but you want to return not more than 480p only due to the conditions of distribution of content to a specific end-user (i.e. free account):\n- To a generic ```.../master.m3u8``` manifest\n- Add a suffix-option to limit quality ```.../master-max-480.m3u8```\n- Add a suffix-option to limit quality and codec ```.../master-min-320-max-320-h264.m3u8```\n  \nCaution. Solely master.m3u8 (and master[-options].m3u8) is officially documented and intended for your use. Any additional internal manifests, sub-manifests, parameters, chunk names, file extensions, and related components are internal infrastructure entities. These may undergo modifications without prior notice, in any manner or form. It is strongly advised not to store them in your database or cache them on your end.",
				Computed:    true,
			},
			"id": schema.Int64Attribute{
				Description: "Video ID",
				Computed:    true,
			},
			"iframe_url": schema.StringAttribute{
				Description: "A URL to a built-in HTML video player with the video inside. It can be inserted into an iframe on your website and the video will automatically play in all browsers.\nThe player can be opened or shared via this direct link. Also the video player can be integrated into your web pages using the Iframe tag.\nExample of usage on a web page:\n<iframe width=\"100%\" height=\"100%\" src=\"https://player.gvideo.co/videos/`2675_FnlHXwA16ZMxmUr`\" frameborder=\"0\" allow=\"autoplay; encrypted-media\" allowfullscreen></iframe>\n  \nThere are some link modificators you can specify and add manually:\n- ?`no_low_latency` – player is forced to use non-low-latency streams HLS MPEG TS, instead of MPEG-DASH CMAF or HLS/LL-HLS CMAF.\n- ?t=(integer) – time to start playback from specified point in the video. Applicable for VOD only.\n- ?`sub_lang`=(language) – force subtitles to specific language (2 letters ISO 639 code of a language).\n- Read more in the Product Documentation.",
				Computed:    true,
			},
			"origin_size": schema.Int64Attribute{
				Description: "Size of original file",
				Computed:    true,
			},
			"origin_video_duration": schema.Int64Attribute{
				Description: "Original video duration in milliseconds",
				Computed:    true,
			},
			"poster_thumb": schema.StringAttribute{
				Description: `Field contains a link to minimized poster image. Original "poster" image is proportionally scaled to a size of 200 pixels in height.`,
				Computed:    true,
			},
			"recording_started_at": schema.StringAttribute{
				Description: "If the video was saved from a stream, then start time of the stream recording is saved here. Format is date time in ISO 8601",
				Computed:    true,
			},
			"screenshot": schema.StringAttribute{
				Description: "A URL to the default screenshot is here. The image is selected from an array of all screenshots based on the “`screenshot_id`” attribute. If you use your own \"poster\", the link to it will be here too.\nOur video player uses this field to display the static image before the video starts playing. As soon as the user hits \"play\" the image will go away.\nIf you use your own external video player, then you can use the value of this field to set the poster/thumbnail in your player.\nExample:\n- `video_js`.poster: ```api.screenshot```\n- clappr.poster: ```api.screenshot```",
				Computed:    true,
			},
			"slug": schema.StringAttribute{
				Description: "A unique alphanumeric identifier used in public URLs to retrieve and view the video. It is unique for each video, generated randomly and set automatically by the system.\nFormat of usage in URL is \\*.../videos/{`client_id`}_{slug}/...\\*\nExample:\n- Player: /videos/`12345_neAq1bYZ2`\n- Manifest: /videos/`12345_neAq1bYZ2`/master.m3u8\n- Rendition: /videos/`12345_neAq1bYZ2`/`qid90v1_720`.mp4",
				Computed:    true,
			},
			"sprite": schema.StringAttribute{
				Description: "Link to picture with video storyboard. Image in JPG format. The picture is a set of rectangles with frames from the video. Typically storyboard is used to show preview images when hovering the video's timeline.",
				Computed:    true,
			},
			"sprite_vtt": schema.StringAttribute{
				Description: "Storyboard in VTT format. This format implies an explicit indication of the timing and frame area from a large sprite image.",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Video processing status:\n- empty – initial status, when video-entity is created, but video-file has not yet been fully uploaded (TUS uploading, or downloading from an origin is not finished yet)\n- pending – video is in queue to be processed\n- viewable – video has at least 1 quality and can already be viewed via a link, but not all qualities are ready yet\n- ready – video is completely ready, available for viewing with all qualities\n- error – error while processing a video, look at \"error\" field\nAvailable values: \"empty\", \"pending\", \"viewable\", \"ready\", \"error\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"empty",
						"pending",
						"viewable",
						"ready",
						"error",
					),
				},
			},
			"stream_id": schema.Int64Attribute{
				Description: "If the video was saved from a stream, then ID of that stream is saved here",
				Computed:    true,
			},
			"views": schema.Int64Attribute{
				Description: "Number of video views through the built-in HTML video player of the Streaming Platform only. This attribute does not count views from other external players and native OS players, so here may be less number of views than in \"`cdn_views`\".",
				Computed:    true,
			},
			"screenshots": schema.ListAttribute{
				Description: "Array of auto generated screenshots from the video. By default 5 static screenshots are taken from different places in the video. If the video is short, there may be fewer screenshots.\nScreenshots are created automatically, so they may contain not very good frames from the video. To use your own image look at \"poster\" attribute.",
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
			"converted_videos": schema.ListNestedAttribute{
				Description: "Array of data about each transcoded quality",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[StreamingVideoConvertedVideosModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "ID of the converted file of the specific quality",
							Computed:    true,
						},
						"error": schema.StringAttribute{
							Description: "Video processing error text in this quality",
							Computed:    true,
						},
						"height": schema.Int64Attribute{
							Description: "Height in pixels of the converted video file of the specific quality. Can be ```null``` for audio-only files.",
							Computed:    true,
						},
						"mp4_url": schema.StringAttribute{
							Description: "A URL to a rendition file of the specified quality in MP4 format for downloading.\n  \n**Download methods**\nFor each converted video, additional download endpoints are available under `converted_videos`/`mp4_urls`.\nAn MP4 download enpoints:\n- /videos/{`client_id`}_{slug}/{filename}.mp4\n- /videos/{`client_id`}_{slug}/{filename}.mp4/download\n- /videos/{`client_id`}_{slug}/{filename}.mp4/download={`custom_filename`}\nThe first option returns the file as is.\nThe following options respond with the header that directly tells browsers to download the file instead of playing it in the browser.\n```\nContent-Disposition: attachment\n```\nThe third option allows you to set a custom name for the file being downloaded. You can optionally specify a custom filename (just name excluding the .mp4 extension) using the download= query.\nFilename Constraints\n- Length: 1-255 characters\n- Must NOT include the .mp4 extension (it is added automatically)\n- Allowed characters: a-z, A-Z, 0-9, _(underscore), -(dash), .(dot)\n- First character cannot be .(dot)\nExample valid filenames: ```holiday2025```, ```_backup.final```, ```clip-v1.2```\n  \n**Default MP4 file name structure**\nLink to the file {filename} contains information about the encoding method using format:\n```___.mp4```\n- `````` – Internal quality identifier and file version. Please do not use it, can be changed at any time without any notice.\n- `````` – Codec name that was used to encode the video, or audio codec if it is an audio-only file.\n- `````` – Encoding bitrate in Kbps.\n- `````` – Video height, or word \"audio\" if it is an audio-only file.\nNote that this link format has been applied since 14.08.2024. If the video entity was uploaded earlier, links may have old simplified format.\nExample: ``` /videos/{`client_id`}_{slug}/`qid3567v1_h264_4050_1080`.mp4 ```\n  \n**Dynamic speed limiting**\nThis mode sets different limits for different users or for different types of content. The speed is adjusted based on requests with the “speed” and “buffer” arguments.\nExample: ``` ?speed=50k&buffer=500k ```\nRead more in Product Documentation in CDN section \"Network limits\".\n  \n**Secure token authentication (updated)**\nAccess to MP4 download links can be protected using secure tokens passed as query parameters. The token generation logic has been updated to allow fine-grained protection per file and bitrate.\nToken generation uses the entire MP4 path, which ensures the token only grants access to a specific quality/version of the video. This prevents unintended access to other bitrate versions of an ABR stream.\nToken Query Parameters:\n- token: The generated hash\n- expires: Expiration timestamp\n- speed: (optional) Speed limit in bytes/sec, or empty string\n- buffer: (optional) Buffer size in bytes, or empty string\nOptional (for IP-bound tokens):\n- ip: The user’s IP address\nExample: ``` ?md5=QX39c77lbQKvYgMMAvpyMQ&expires=1743167062 ```\nRead more in Product Documentation in Streaming section \"Protected temporarily link\".\n  \n**Examples**\n- Audio-only: ```https://demo-public.gvideo.io/videos/`2675_JNnccG5l97XPxsov`/`qid3585v1_aac_128_audio`.mp4```\n- Video: ```https://demo-public.gvideo.io/videos/`2675_3MlggU4xDb1Ssa5Y`/`qid3567v1_h264_4050_1080`.mp4/download```\n- Video with custom download filename: ```https://demo-public.gvideo.io/videos/`2675_XtMKxzJM6Xt7SBUV`/1080.mp4/download=`highlights_v1`.`1_2025`-05-30```",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Specific quality name",
							Computed:    true,
						},
						"progress": schema.Int64Attribute{
							Description: "Status of transcoding into the specific quality, from 0 to 100",
							Computed:    true,
						},
						"size": schema.Int64Attribute{
							Description: "Size in bytes of the converted file of the specific quality. Can be ```null``` until transcoding is fully completed.",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "Status of transcoding:\n- processing – video is being transcoded to this quality,\n- complete – quality is fully processed,\n- error – quality processing error, see parameter \"error\".\nAvailable values: \"processing\", \"complete\", \"error\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"processing",
									"complete",
									"error",
								),
							},
						},
						"width": schema.Int64Attribute{
							Description: "Width in pixels of the converted video file of the specified quality. Can be ```null``` for audio files.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (r *StreamingVideoResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *StreamingVideoResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
