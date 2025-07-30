// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_template

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ resource.ResourceWithConfigValidators = (*FastedgeTemplateResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:   "Template ID",
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown(), int64planmodifier.RequiresReplace()},
			},
			"binary_id": schema.Int64Attribute{
				Description:   "Binary ID",
				Required:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Description:   "Name of the template",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"owned": schema.BoolAttribute{
				Description:   "Is the template owned by user?",
				Required:      true,
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplace()},
			},
			"params": schema.ListNestedAttribute{
				Description: "Parameters",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"data_type": schema.StringAttribute{
							Description: "Parameter type\nAvailable values: \"string\", \"number\", \"date\", \"time\", \"secret\".",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"string",
									"number",
									"date",
									"time",
									"secret",
								),
							},
						},
						"mandatory": schema.BoolAttribute{
							Description: "Is this field mandatory?",
							Computed:    true,
							Optional:    true,
							Default:     booldefault.StaticBool(false),
						},
						"name": schema.StringAttribute{
							Description: "Parameter name",
							Required:    true,
						},
						"descr": schema.StringAttribute{
							Description: "Parameter description",
							Optional:    true,
						},
					},
				},
				PlanModifiers: []planmodifier.List{listplanmodifier.RequiresReplace()},
			},
			"long_descr": schema.StringAttribute{
				Description:   "Long description of the template",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"short_descr": schema.StringAttribute{
				Description:   "Short description of the template",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"api_type": schema.StringAttribute{
				Description: "Wasm API type",
				Computed:    true,
			},
		},
	}
}

func (r *FastedgeTemplateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *FastedgeTemplateResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
