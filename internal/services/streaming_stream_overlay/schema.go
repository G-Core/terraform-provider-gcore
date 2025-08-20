// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_stream_overlay

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ resource.ResourceWithConfigValidators = (*StreamingStreamOverlayResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"stream_id": schema.Int64Attribute{
				Required:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"overlay_id": schema.Int64Attribute{
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"body": schema.ListNestedAttribute{
				Computed:   true,
				Optional:   true,
				CustomType: customfield.NewNestedObjectListType[StreamingStreamOverlayBodyModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"url": schema.StringAttribute{
							Description: "Valid http/https URL to an HTML page/widget",
							Required:    true,
						},
						"height": schema.Int64Attribute{
							Description: "Height of the widget",
							Optional:    true,
						},
						"stretch": schema.BoolAttribute{
							Description: `Switch of auto scaling the widget. Must not be used as "true" simultaneously with the coordinate installation method (w, h, x, y).`,
							Computed:    true,
							Optional:    true,
							Default:     booldefault.StaticBool(false),
						},
						"width": schema.Int64Attribute{
							Description: "Width of the widget",
							Optional:    true,
						},
						"x": schema.Int64Attribute{
							Description: "Coordinate of left upper corner",
							Optional:    true,
						},
						"y": schema.Int64Attribute{
							Description: "Coordinate of left upper corner",
							Optional:    true,
						},
					},
				},
				PlanModifiers: []planmodifier.List{listplanmodifier.RequiresReplaceIfConfigured()},
			},
			"height": schema.Int64Attribute{
				Description: "Height of the widget",
				Optional:    true,
			},
			"url": schema.StringAttribute{
				Description: "Valid http/https URL to an HTML page/widget",
				Optional:    true,
			},
			"width": schema.Int64Attribute{
				Description: "Width of the widget",
				Optional:    true,
			},
			"x": schema.Int64Attribute{
				Description: "Coordinate of left upper corner",
				Optional:    true,
			},
			"y": schema.Int64Attribute{
				Description: "Coordinate of left upper corner",
				Optional:    true,
			},
			"stretch": schema.BoolAttribute{
				Description: `Switch of auto scaling the widget. Must not be used as "true" simultaneously with the coordinate installation method (w, h, x, y).`,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"created_at": schema.StringAttribute{
				Description: "Datetime of creation in ISO 8601",
				Computed:    true,
			},
			"id": schema.Int64Attribute{
				Description: "ID of the overlay",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Datetime of last update in ISO 8601",
				Computed:    true,
			},
		},
	}
}

func (r *StreamingStreamOverlayResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *StreamingStreamOverlayResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
