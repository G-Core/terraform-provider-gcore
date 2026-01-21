// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_insight_silence

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigValidators = (*WaapDomainInsightSilenceResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:   "A generated unique identifier for the silence",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"domain_id": schema.Int64Attribute{
				Description:   "The domain ID",
				Required:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"insight_type": schema.StringAttribute{
				Description:   "The slug of the insight type",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"author": schema.StringAttribute{
				Description: "The author of the silence",
				Required:    true,
			},
			"comment": schema.StringAttribute{
				Description: "A comment explaining the reason for the silence",
				Required:    true,
			},
			"labels": schema.MapAttribute{
				Description: "A hash table of label names and values that apply to the insight silence",
				Required:    true,
				ElementType: types.StringType,
			},
			"expire_at": schema.StringAttribute{
				Description: "The date and time the silence expires in ISO 8601 format",
				Optional:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
		},
	}
}

func (r *WaapDomainInsightSilenceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *WaapDomainInsightSilenceResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
