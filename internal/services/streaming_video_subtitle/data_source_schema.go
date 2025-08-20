// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_video_subtitle

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSourceWithConfigValidators = (*StreamingVideoSubtitleDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Required: true,
			},
			"video_id": schema.Int64Attribute{
				Required: true,
			},
			"language": schema.StringAttribute{
				Description: "3-letter language code according to ISO-639-2 (bibliographic code)",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of subtitle file",
				Computed:    true,
			},
			"vtt": schema.StringAttribute{
				Description: `Full text of subtitles/captions, with escaped "\n" ("\r") symbol of new line`,
				Computed:    true,
			},
		},
	}
}

func (d *StreamingVideoSubtitleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *StreamingVideoSubtitleDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
