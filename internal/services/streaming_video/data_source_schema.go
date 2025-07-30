// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_video

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*StreamingVideoDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"video_id": schema.Int64Attribute{
				Required: true,
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
			"client_user_id": schema.Int64Attribute{
				Description: "Custom meta field for storing the Identifier in your system. We do not use this field in any way when processing the stream. Example: ```client_user_id = 1001```",
				Computed:    true,
			},
			"custom_iframe_url": schema.StringAttribute{
				Description: "Custom URL of Iframe for video player to be used in share panel in player. Auto generated Iframe URL provided by default.",
				Computed:    true,
			},
			"dash_url": schema.StringAttribute{
				Description: "A URL to a master playlist MPEG-DASH (master.mpd) with CMAF or WebM based chunks.\nChunk type will be selected automatically for each quality:\n- CMAF for H264 and H265 codecs.\n- WebM for AV1 codec.\n  \nThis URL is a link to the main manifest. But you can also manually specify suffix-options that will allow you to change the manifest to your request:\n```/videos/{`client_id`}_{slug}/master[-min-N][-max-N][-(h264|hevc|av1)].mpd```\nList of suffix-options:\n- [-min-N] – ABR soft limitation of qualities from below.\n- [-max-N] – ABR soft limitation of qualities from above.\n- [-(h264|hevc|av1) – Video codec soft limitation. Applicable if the video was transcoded into multiple codecs H264, H265 and AV1 at once, but you want to return just 1 video codec in a manifest. Read the Product Documentation for details.\nRead more what is ABR soft-limiting in the \"`hls_url`\" field above.\n  \nCaution. Solely master.mpd is officially documented and intended for your use. Any additional internal manifests, sub-manifests, parameters, chunk names, file extensions, and related components are internal infrastructure entities. These may undergo modifications without prior notice, in any manner or form. It is strongly advised not to store them in your database or cache them on your end.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Additional text field for video description",
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
			"name": schema.StringAttribute{
				Description: "Title of the video.\nOften used as a human-readable name of the video, but can contain any text you wish. The values are not unique and may be repeated.\nExamples:\n- Educational training 2024-03-29\n- Series X S3E14, The empire strikes back\n- 480fd499-2de2-4988-bc1a-a4eebe9818ee",
				Computed:    true,
			},
			"origin_size": schema.Int64Attribute{
				Description: "Size of original file",
				Computed:    true,
			},
			"origin_url": schema.StringAttribute{
				Description: "URL to an original file from which the information for transcoding was taken.\nMay contain a link for scenarios:\n- If the video was downloaded from another origin\n- If the video is a recording of a live stream\n- Otherwise it is \"null\"\n**Copy from another server**\nURL to an original file that was downloaded. Look at method \"Copy from another server\" in POST /videos.\n**Recording of an original live stream**\nURL to the original non-transcoded stream recording with original quality, saved in MP4 format. File is created immediately after the completion of the stream recording. The stream from which the recording was made is reflected in \"`stream_id`\" field.\nCan be used for internal operations when a recording needs to be received faster than the transcoded versions are ready. But this version is not intended for public distribution. Views and downloads occur in the usual way, like viewing an MP4 rendition.\nThe MP4 file becomes available for downloading when the video entity \"status\" changes from \"new\" to \"pending\". The file is stored for 7 days, after which it will be automatically deleted.\nFormat of URL is ```/videos/_/`origin__`.mp4```\nWhere:\n- `````` – Encoding bitrate in Kbps.\n- `````` – Video height.\nThis is a premium feature, available only upon request through your manager or support team.",
				Computed:    true,
			},
			"origin_video_duration": schema.Int64Attribute{
				Description: "Original video duration in milliseconds",
				Computed:    true,
			},
			"poster": schema.StringAttribute{
				Description: "Poster is your own static image which can be displayed before the video begins playing. This is often a frame of the video or a custom title screen.\nField contains a link to your own uploaded image.\nAlso look at \"screenshot\" attribute.",
				Computed:    true,
			},
			"poster_thumb": schema.StringAttribute{
				Description: `Field contains a link to minimized poster image. Original "poster" image is proportionally scaled to a size of 200 pixels in height.`,
				Computed:    true,
			},
			"projection": schema.StringAttribute{
				Description: "Regulates the video format:\n\n* **regular** — plays the video as usual\n* **vr360** — plays the video in 360 degree mode\n* **vr180** — plays the video in 180 degree mode\n* **vr360tb** — plays the video in 3D 360 degree mode Top-Bottom.\n\n  \n Default is regular",
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
			"screenshot_id": schema.Int64Attribute{
				Description: "ID of auto generated screenshots to be used for default screenshot.\nCounting from 0. A value of -1 sets the \"screenshot\" attribute to the URL of your own image from the \"poster\" attribute.",
				Computed:    true,
			},
			"share_url": schema.StringAttribute{
				Description: "Custom URL or iframe displayed in the link field when a user clicks on a sharing button in player. If empty, the link field and social network sharing is disabled",
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
				CustomType:  customfield.NewNestedObjectListType[StreamingVideoConvertedVideosDataSourceModel](ctx),
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

func (d *StreamingVideoDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *StreamingVideoDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
