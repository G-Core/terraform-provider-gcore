// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_app

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigValidators = (*FastedgeAppResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "Apps are descriptions of edge apps, that reference the binary and may contain app-specific settings, such as environment variables.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:   "App ID",
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"template": schema.Int64Attribute{
				Description: "Template ID",
				Optional:    true,
			},
			"binary": schema.Int64Attribute{
				Description: "Binary ID",
				Computed:    true,
				Optional:    true,
			},
			"comment": schema.StringAttribute{
				Description: "App description",
				Computed:    true,
				Optional:    true,
			},
			"debug": schema.BoolAttribute{
				Description: "Switch on logging for 30 minutes (switched off by default)",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"log": schema.StringAttribute{
				Description: "Logging channel (by default - kafka, which allows exploring logs with API)\nAvailable values: \"kafka\", \"none\".",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("kafka", "none"),
				},
			},
			"name": schema.StringAttribute{
				Description: "App name",
				Computed:    true,
				Optional:    true,
			},
			"status": schema.Int64Attribute{
				Description: "Status code:  \n0 - draft (inactive)  \n1 - enabled  \n2 - disabled  \n3 - hourly call limit exceeded  \n4 - daily call limit exceeded  \n5 - suspended",
				Computed:    true,
				Optional:    true,
			},
			"env": schema.MapAttribute{
				Description: "Environment variables",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewMapType[types.String](ctx),
				ElementType: types.StringType,
			},
			"rsp_headers": schema.MapAttribute{
				Description: "Extra headers to add to the response",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewMapType[types.String](ctx),
				ElementType: types.StringType,
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
			"stores": schema.MapNestedAttribute{
				Description: "Application edge stores",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewNestedObjectMapType[FastedgeAppStoresModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "The identifier of the store",
							Required:    true,
						},
						"comment": schema.StringAttribute{
							Description: "A description of the store",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the store",
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
