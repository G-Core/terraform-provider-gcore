// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_playlist

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*StreamingPlaylistsDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"max_items": schema.Int64Attribute{
				Description: "Max items to fetch, default: 1000",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"items": schema.ListNestedAttribute{
				Description: "The items returned by the data source",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[StreamingPlaylistsItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"active": schema.BoolAttribute{
							Description: "Enables/Disables playlist. Has two possible values:\n- true – Playlist can be played.\n- false – Playlist is disabled. No broadcast while it's desabled.",
							Computed:    true,
						},
						"ad_id": schema.Int64Attribute{
							Description: "The advertisement ID that will be inserted into the video",
							Computed:    true,
						},
						"client_id": schema.Int64Attribute{
							Description: "Current playlist client ID",
							Computed:    true,
						},
						"client_user_id": schema.Int64Attribute{
							Description: "Custom field where you can specify user ID in your system",
							Computed:    true,
						},
						"countdown": schema.BoolAttribute{
							Description: "Enables countdown before playlist start with ```playlist_type: live```",
							Computed:    true,
						},
						"hls_cmaf_url": schema.StringAttribute{
							Description: "A URL to a master playlist HLS (master-cmaf.m3u8) with CMAF-based chunks. Chunks are in fMP4 container.\n  \nIt is possible to use the same suffix-options as described in the \"`hls_url`\" attribute.\n  \nCaution. Solely master.m3u8 (and master[-options].m3u8) is officially documented and intended for your use. Any additional internal manifests, sub-manifests, parameters, chunk names, file extensions, and related components are internal infrastructure entities. These may undergo modifications without prior notice, in any manner or form. It is strongly advised not to store them in your database or cache them on your end.",
							Computed:    true,
						},
						"hls_url": schema.StringAttribute{
							Description: "A URL to a master playlist HLS (master.m3u8) with MPEG TS container.\n  \nThis URL is a link to the main manifest. But you can also manually specify suffix-options that will allow you to change the manifest to your request:\n```/playlists/{`client_id`}_{`playlist_id`}/master[-cmaf][-min-N][-max-N][-img][-(h264|hevc|av1)].m3u8```\nPlease see the details in ```hls_url``` attribute of /videos/{id} method.\n  \nCaution. Solely master.m3u8 (and master[-options].m3u8) is officially documented and intended for your use. Any additional internal manifests, sub-manifests, parameters, chunk names, file extensions, and related components are internal infrastructure entities. These may undergo modifications without prior notice, in any manner or form. It is strongly advised not to store them in your database or cache them on your end.",
							Computed:    true,
						},
						"iframe_url": schema.StringAttribute{
							Description: "A URL to a built-in HTML video player with the video inside. It can be inserted into an iframe on your website and the video will automatically play in all browsers.\nThe player can be opened or shared via this direct link. Also the video player can be integrated into your web pages using the Iframe tag.\n  \nPlease see the details in ```iframe_url``` attribute of /videos/{id} method.",
							Computed:    true,
						},
						"loop": schema.BoolAttribute{
							Description: "Enables/Disables playlist loop",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Playlist name",
							Computed:    true,
						},
						"player_id": schema.Int64Attribute{
							Description: "The player ID with which the video will be played",
							Computed:    true,
						},
						"playlist_type": schema.StringAttribute{
							Description: "Determines whether the playlist:\n- `live` - playlist for live-streaming\n- `vod` - playlist is for video on demand access\nAvailable values: \"live\", \"vod\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("live", "vod"),
							},
						},
						"start_time": schema.StringAttribute{
							Description: "Playlist start time. Playlist won't be available before the specified time. Datetime in ISO 8601 format.",
							Computed:    true,
						},
						"video_ids": schema.ListAttribute{
							Description: "A list of VOD IDs included in the playlist. Order of videos in a\nplaylist reflects the order of IDs in the array.\nMaximum video limit = 128.",
							Computed:    true,
							CustomType:  customfield.NewListType[types.Int64](ctx),
							ElementType: types.Int64Type,
						},
					},
				},
			},
		},
	}
}

func (d *StreamingPlaylistsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *StreamingPlaylistsDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
