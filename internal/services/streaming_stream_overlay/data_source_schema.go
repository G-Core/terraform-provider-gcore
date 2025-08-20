// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_stream_overlay

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSourceWithConfigValidators = (*StreamingStreamOverlayDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"overlay_id": schema.Int64Attribute{
				Required: true,
			},
			"stream_id": schema.Int64Attribute{
				Required: true,
			},
			"created_at": schema.StringAttribute{
				Description: "Datetime of creation in ISO 8601",
				Computed:    true,
			},
			"height": schema.Int64Attribute{
				Description: "Height of the widget",
				Computed:    true,
			},
			"id": schema.Int64Attribute{
				Description: "ID of the overlay",
				Computed:    true,
			},
			"stretch": schema.BoolAttribute{
				Description: `Switch of auto scaling the widget. Must not be used as "true" simultaneously with the coordinate installation method (w, h, x, y).`,
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
	}
}

func (d *StreamingStreamOverlayDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *StreamingStreamOverlayDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
