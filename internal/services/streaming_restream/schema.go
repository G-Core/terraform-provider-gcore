// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_restream

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

var _ resource.ResourceWithConfigValidators = (*StreamingRestreamResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"restream_id": schema.Int64Attribute{
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"restream": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"active": schema.BoolAttribute{
						Description: "Enables/Disables restream. Has two possible values:\n\n* **true** — restream is enabled and can be started\n* **false** — restream is disabled.\n\n  \nDefault is true",
						Optional:    true,
					},
					"client_user_id": schema.Int64Attribute{
						Description: "Custom field where you can specify user ID in your system",
						Optional:    true,
					},
					"live": schema.BoolAttribute{
						Description: "Indicates that the stream is being published. Has two possible values:\n\n* **true** — stream is being published\n* **false** — stream isn't published",
						Optional:    true,
					},
					"name": schema.StringAttribute{
						Description: "Restream name",
						Optional:    true,
					},
					"stream_id": schema.Int64Attribute{
						Description: "ID of the stream to restream",
						Optional:    true,
					},
					"uri": schema.StringAttribute{
						Description: "A URL to push the stream to",
						Optional:    true,
					},
				},
			},
			"active": schema.BoolAttribute{
				Description: "Enables/Disables restream. Has two possible values:\n\n* **true** — restream is enabled and can be started\n* **false** — restream is disabled.\n\n  \nDefault is true",
				Computed:    true,
			},
			"client_user_id": schema.Int64Attribute{
				Description: "Custom field where you can specify user ID in your system",
				Computed:    true,
			},
			"live": schema.BoolAttribute{
				Description: "Indicates that the stream is being published. Has two possible values:\n\n* **true** — stream is being published\n* **false** — stream isn't published",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Restream name",
				Computed:    true,
			},
			"stream_id": schema.Int64Attribute{
				Description: "ID of the stream to restream",
				Computed:    true,
			},
			"uri": schema.StringAttribute{
				Description: "A URL to push the stream to",
				Computed:    true,
			},
		},
	}
}

func (r *StreamingRestreamResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *StreamingRestreamResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
