package cdn_origin_shielding

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Origin shielding protects your origin server from being overloaded by routing all CDN requests through a single shield (precache) server for a CDN resource.\n\nOrigin shielding is a paid option. Accounts without it enabled receive an API error when enabling shielding.",
		Attributes: map[string]schema.Attribute{
			"resource_id": schema.Int64Attribute{
				Required:      true,
				Description:   "ID of the CDN resource for which origin shielding is configured. Changing this forces a new resource.",
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"shielding_pop": schema.Int64Attribute{
				Required:    true,
				Description: "Origin shielding location ID (point of presence). Look up available IDs with the gcore_cdn_origin_shielding data source or the GET /cdn/shieldingpop_v2 endpoint.",
			},
		},
	}
}

func (r *CDNOriginShieldingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}
