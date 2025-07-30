// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package iam_api_token

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ resource.ResourceWithConfigValidators = (*IamAPITokenResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:   "API token.\nCopy it, because you will not be able to get it again.\nWe do not store tokens. All responsibility for token storage and usage is on the issuer.",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown(), stringplanmodifier.RequiresReplace()},
			},
			"token": schema.StringAttribute{
				Description:   "API token.\nCopy it, because you will not be able to get it again.\nWe do not store tokens. All responsibility for token storage and usage is on the issuer.",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown(), stringplanmodifier.RequiresReplace()},
			},
			"client_id": schema.Int64Attribute{
				Required:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"token_id": schema.Int64Attribute{
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"exp_date": schema.StringAttribute{
				Description:   "Date when the API token becomes expired (ISO 8086/RFC 3339 format), UTC.\nIf null, then the API token will never expire.",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Description:   "API token name.",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"client_user": schema.SingleNestedAttribute{
				Description: "API token role.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"role": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"id": schema.Int64Attribute{
								Description: "Group's ID: Possible values are:   \n\n* 1 - Administrators* 2 - Users* 5 - Engineers* 3009 - Purge and Prefetch only (API+Web)* 3022 - Purge and Prefetch only (API)",
								Optional:    true,
							},
							"name": schema.StringAttribute{
								Description: "Group's name.\nAvailable values: \"Users\", \"Administrators\", \"Engineers\", \"Purge and Prefetch only (API)\", \"Purge and Prefetch only (API+Web)\".",
								Optional:    true,
								Validators: []validator.String{
									stringvalidator.OneOfCaseInsensitive(
										"Users",
										"Administrators",
										"Engineers",
										"Purge and Prefetch only (API)",
										"Purge and Prefetch only (API+Web)",
									),
								},
							},
						},
					},
				},
				PlanModifiers: []planmodifier.Object{objectplanmodifier.RequiresReplace()},
			},
			"description": schema.StringAttribute{
				Description:   "API token description.",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"created": schema.StringAttribute{
				Description: "Date when the API token was issued (ISO 8086/RFC 3339 format), UTC.",
				Computed:    true,
			},
			"deleted": schema.BoolAttribute{
				Description: "Deletion flag. If true, then the API token was deleted.",
				Computed:    true,
			},
			"expired": schema.BoolAttribute{
				Description: "Expiration flag. If true, then the API token has expired.\nWhen an API token expires it will be automatically deleted.",
				Computed:    true,
			},
			"last_usage": schema.StringAttribute{
				Description: "Date when the API token was last used (ISO 8086/RFC 3339 format), UTC.",
				Computed:    true,
			},
		},
	}
}

func (r *IamAPITokenResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *IamAPITokenResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
