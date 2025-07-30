// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package iam_api_token

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*IamAPITokenDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"client_id": schema.Int64Attribute{
				Required: true,
			},
			"token_id": schema.Int64Attribute{
				Required: true,
			},
			"created": schema.StringAttribute{
				Description: "Date when the API token was issued (ISO 8086/RFC 3339 format), UTC.",
				Computed:    true,
			},
			"deleted": schema.BoolAttribute{
				Description: "Deletion flag. If true, then the API token was deleted.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "API token description.",
				Computed:    true,
			},
			"exp_date": schema.StringAttribute{
				Description: "Date when the API token becomes expired (ISO 8086/RFC 3339 format), UTC.\nIf null, then the API token will never expire.",
				Computed:    true,
			},
			"expired": schema.BoolAttribute{
				Description: "Expiration flag. If true, then the API token has expired.\nWhen an API token expires it will be automatically deleted.",
				Computed:    true,
			},
			"id": schema.Int64Attribute{
				Description: "API token ID.",
				Computed:    true,
			},
			"last_usage": schema.StringAttribute{
				Description: "Date when the API token was last used (ISO 8086/RFC 3339 format), UTC.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "API token name.",
				Computed:    true,
			},
			"client_user": schema.SingleNestedAttribute{
				Computed:   true,
				CustomType: customfield.NewNestedObjectType[IamAPITokenClientUserDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"client_id": schema.Int64Attribute{
						Description: "Account's ID.",
						Computed:    true,
					},
					"deleted": schema.BoolAttribute{
						Description: "Deletion flag. If true, then the API token was deleted.",
						Computed:    true,
					},
					"role": schema.SingleNestedAttribute{
						Computed:   true,
						CustomType: customfield.NewNestedObjectType[IamAPITokenClientUserRoleDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"id": schema.Int64Attribute{
								Description: "Group's ID: Possible values are:   \n\n* 1 - Administrators* 2 - Users* 5 - Engineers* 3009 - Purge and Prefetch only (API+Web)* 3022 - Purge and Prefetch only (API)",
								Computed:    true,
							},
							"name": schema.StringAttribute{
								Description: "Group's name.\nAvailable values: \"Users\", \"Administrators\", \"Engineers\", \"Purge and Prefetch only (API)\", \"Purge and Prefetch only (API+Web)\".",
								Computed:    true,
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
					"user_email": schema.StringAttribute{
						Description: "User's email who issued the API token.",
						Computed:    true,
					},
					"user_id": schema.Int64Attribute{
						Description: "User's ID who issued the API token.",
						Computed:    true,
					},
					"user_name": schema.StringAttribute{
						Description: "User's name who issued the API token.",
						Computed:    true,
					},
				},
			},
		},
	}
}

func (d *IamAPITokenDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *IamAPITokenDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
