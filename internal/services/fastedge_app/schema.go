// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_app

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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
		MarkdownDescription: "FastEdge applications combine a WebAssembly binary with configuration, environment variables, and secrets for deployment at the CDN edge.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "App ID",
				Computed:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseNonNullStateForUnknown()},
			},
			"template": schema.Int64Attribute{
				Description: "Template ID",
				Optional:    true,
			},
			"binary": schema.Int64Attribute{
				Description: "ID of the WebAssembly binary to deploy",
				Computed:    true,
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"comment": schema.StringAttribute{
				Description: "Optional human-readable description of the application's purpose",
				Computed:    true,
				Optional:    true,
			},
			"debug": schema.BoolAttribute{
				Description: "Enable verbose debug logging for 30 minutes. Automatically expires to prevent performance impact.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"name": schema.StringAttribute{
				Description: "Unique application name (alphanumeric, hyphens allowed)",
				Computed:    true,
				Optional:    true,
			},
			"status": schema.Int64Attribute{
				Description: "Status code:  \n0 - draft (inactive)  \n1 - enabled  \n2 - disabled  \n5 - suspended",
				Computed:    true,
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(0, 5),
				},
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
					},
				},
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
