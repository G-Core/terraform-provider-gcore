// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_app

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ resource.ResourceWithConfigValidators = (*FastedgeAppResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:   "App ID",
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"binary": schema.Int64Attribute{
				Description: "Binary ID",
				Optional:    true,
			},
			"comment": schema.StringAttribute{
				Description: "App description",
				Optional:    true,
			},
			"log": schema.StringAttribute{
				Description: "Logging channel (by default - kafka, which allows exploring logs with API)\nAvailable values: \"kafka\", \"none\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("kafka", "none"),
				},
			},
			"name": schema.StringAttribute{
				Description: "App name",
				Optional:    true,
			},
			"status": schema.Int64Attribute{
				Description: "Status code:  \n0 - draft (inactive)  \n1 - enabled  \n2 - disabled  \n3 - hourly call limit exceeded  \n4 - daily call limit exceeded  \n5 - suspended",
				Optional:    true,
			},
			"template": schema.Int64Attribute{
				Description: "Template ID",
				Optional:    true,
			},
			"env": schema.MapAttribute{
				Description: "Environment variables",
				Optional:    true,
				ElementType: types.StringType,
			},
			"rsp_headers": schema.MapAttribute{
				Description: "Extra headers to add to the response",
				Optional:    true,
				ElementType: types.StringType,
			},
			"stores": schema.MapAttribute{
				Description: "KV stores for the app",
				Optional:    true,
				ElementType: types.Int64Type,
			},
			"debug": schema.BoolAttribute{
				Description: "Switch on logging for 30 minutes (switched off by default)",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"secrets": schema.MapNestedAttribute{
				Description: "Application secrets",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewNestedObjectMapType[FastedgeAppSecretsModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "The unique identifier of the secret.",
							Required:    true,
						},
						"comment": schema.StringAttribute{
							Description: "A description or comment about the secret.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The unique name of the secret.",
							Computed:    true,
						},
					},
				},
			},
			"api_type": schema.StringAttribute{
				Description: "Wasm API type",
				Computed:    true,
			},
			"debug_until": schema.StringAttribute{
				Description: "When debugging finishes",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"plan": schema.StringAttribute{
				Description: "Plan name",
				Computed:    true,
			},
			"plan_id": schema.Int64Attribute{
				Description: "Plan ID",
				Computed:    true,
			},
			"template_name": schema.StringAttribute{
				Description: "Template name",
				Computed:    true,
			},
			"upgradeable_to": schema.Int64Attribute{
				Description: "ID of the binary the app can be upgraded to",
				Computed:    true,
			},
			"url": schema.StringAttribute{
				Description: "App URL",
				Computed:    true,
			},
			"networks": schema.ListAttribute{
				Description: "Networks",
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
		},
	}
}

func (r *FastedgeAppResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *FastedgeAppResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
