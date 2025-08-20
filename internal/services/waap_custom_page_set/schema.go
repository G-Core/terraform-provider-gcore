// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_custom_page_set

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigValidators = (*WaapCustomPageSetResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:   "The ID of the custom page set",
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Description: "Name of the custom page set",
				Required:    true,
			},
			"domains": schema.ListAttribute{
				Description: "List of domain IDs that are associated with this page set",
				Optional:    true,
				ElementType: types.Int64Type,
			},
			"block": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description: "Indicates whether the custom custom page is active or inactive",
						Required:    true,
					},
					"header": schema.StringAttribute{
						Description: "The text to display in the header of the custom page",
						Optional:    true,
					},
					"logo": schema.StringAttribute{
						Description: "Supported image types are JPEG, PNG and JPG, size is limited to width 450px, height 130px. This should be a base 64 encoding of the full HTML img tag compatible image, with the header included.",
						Optional:    true,
					},
					"text": schema.StringAttribute{
						Description: "The text to display in the body of the custom page",
						Optional:    true,
					},
					"title": schema.StringAttribute{
						Description: "The text to display in the title of the custom page",
						Optional:    true,
					},
				},
			},
			"block_csrf": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description: "Indicates whether the custom custom page is active or inactive",
						Required:    true,
					},
					"header": schema.StringAttribute{
						Description: "The text to display in the header of the custom page",
						Optional:    true,
					},
					"logo": schema.StringAttribute{
						Description: "Supported image types are JPEG, PNG and JPG, size is limited to width 450px, height 130px. This should be a base 64 encoding of the full HTML img tag compatible image, with the header included.",
						Optional:    true,
					},
					"text": schema.StringAttribute{
						Description: "The text to display in the body of the custom page",
						Optional:    true,
					},
					"title": schema.StringAttribute{
						Description: "The text to display in the title of the custom page",
						Optional:    true,
					},
				},
			},
			"captcha": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description: "Indicates whether the custom custom page is active or inactive",
						Required:    true,
					},
					"error": schema.StringAttribute{
						Description: "Error message",
						Optional:    true,
					},
					"header": schema.StringAttribute{
						Description: "The text to display in the header of the custom page",
						Optional:    true,
					},
					"logo": schema.StringAttribute{
						Description: "Supported image types are JPEG, PNG and JPG, size is limited to width 450px, height 130px. This should be a base 64 encoding of the full HTML img tag compatible image, with the header included.",
						Optional:    true,
					},
					"text": schema.StringAttribute{
						Description: "The text to display in the body of the custom page",
						Optional:    true,
					},
					"title": schema.StringAttribute{
						Description: "The text to display in the title of the custom page",
						Optional:    true,
					},
				},
			},
			"cookie_disabled": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description: "Indicates whether the custom custom page is active or inactive",
						Required:    true,
					},
					"header": schema.StringAttribute{
						Description: "The text to display in the header of the custom page",
						Optional:    true,
					},
					"text": schema.StringAttribute{
						Description: "The text to display in the body of the custom page",
						Optional:    true,
					},
				},
			},
			"handshake": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description: "Indicates whether the custom custom page is active or inactive",
						Required:    true,
					},
					"header": schema.StringAttribute{
						Description: "The text to display in the header of the custom page",
						Optional:    true,
					},
					"logo": schema.StringAttribute{
						Description: "Supported image types are JPEG, PNG and JPG, size is limited to width 450px, height 130px. This should be a base 64 encoding of the full HTML img tag compatible image, with the header included.",
						Optional:    true,
					},
					"title": schema.StringAttribute{
						Description: "The text to display in the title of the custom page",
						Optional:    true,
					},
				},
			},
			"javascript_disabled": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description: "Indicates whether the custom custom page is active or inactive",
						Required:    true,
					},
					"header": schema.StringAttribute{
						Description: "The text to display in the header of the custom page",
						Optional:    true,
					},
					"text": schema.StringAttribute{
						Description: "The text to display in the body of the custom page",
						Optional:    true,
					},
				},
			},
		},
	}
}

func (r *WaapCustomPageSetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *WaapCustomPageSetResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
