// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_kv_store

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ resource.ResourceWithConfigValidators = (*FastedgeKvStoreResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:   "The unique identifier of the store",
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown(), int64planmodifier.RequiresReplace()},
			},
			"comment": schema.StringAttribute{
				Description:   "A description of the store",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"byod": schema.SingleNestedAttribute{
				Description: "BYOD (Bring Your Own Data) settings",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"prefix": schema.StringAttribute{
						Description: "Key prefix",
						Required:    true,
					},
					"url": schema.StringAttribute{
						Description: "URL to connect to",
						Required:    true,
					},
				},
				PlanModifiers: []planmodifier.Object{objectplanmodifier.RequiresReplace()},
			},
			"app_count": schema.Int64Attribute{
				Description: "The number of applications that use this store",
				Computed:    true,
			},
			"updated": schema.StringAttribute{
				Description: "Last update time",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"stats": schema.SingleNestedAttribute{
				Description: "Store statistics",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[FastedgeKvStoreStatsModel](ctx),
				Attributes: map[string]schema.Attribute{
					"cf_count": schema.Int64Attribute{
						Description: "Total number of Cuckoo filter entries",
						Computed:    true,
					},
					"kv_count": schema.Int64Attribute{
						Description: "Total number of KV entries",
						Computed:    true,
					},
					"size": schema.Int64Attribute{
						Description: "Total store size in bytes",
						Computed:    true,
					},
					"zset_count": schema.Int64Attribute{
						Description: "Total number of sorted set entries",
						Computed:    true,
					},
				},
			},
		},
	}
}

func (r *FastedgeKvStoreResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *FastedgeKvStoreResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
