// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package iam_user

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*IamUsersDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"limit": schema.Int64Attribute{
				Description: "The maximum number of items.",
				Optional:    true,
			},
			"max_items": schema.Int64Attribute{
				Description: "Max items to fetch, default: 1000",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"items": schema.ListNestedAttribute{
				Description: "The items returned by the data source",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[IamUsersItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "User's ID.",
							Computed:    true,
						},
						"activated": schema.BoolAttribute{
							Description: "Email confirmation:\n- `true` – user confirmed the email;\n- `false` – user did not confirm the email.",
							Computed:    true,
						},
						"auth_types": schema.ListAttribute{
							Description: "System field. List of auth types available for the account.",
							Computed:    true,
							Validators: []validator.List{
								listvalidator.ValueStringsAre(
									stringvalidator.OneOfCaseInsensitive(
										"password",
										"sso",
										"github",
										"google-oauth2",
									),
								),
							},
							CustomType:  customfield.NewListType[types.String](ctx),
							ElementType: types.StringType,
						},
						"client": schema.Float64Attribute{
							Description: "User's account ID.",
							Computed:    true,
						},
						"company": schema.StringAttribute{
							Description: "User's company.",
							Computed:    true,
						},
						"deleted": schema.BoolAttribute{
							Description: "Deletion flag. If `true` then user was deleted.",
							Computed:    true,
						},
						"email": schema.StringAttribute{
							Description: "User's email address.",
							Computed:    true,
						},
						"groups": schema.ListNestedAttribute{
							Description: "User's group in the current account.\nIAM supports 5 groups:\n- Users\n- Administrators\n- Engineers\n- Purge and Prefetch only (API)\n- Purge and Prefetch only (API+Web)",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectListType[IamUsersGroupsDataSourceModel](ctx),
							NestedObject: schema.NestedAttributeObject{
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
						},
						"lang": schema.StringAttribute{
							Description: "User's language.\nDefines language of the control panel and email messages.\nAvailable values: \"de\", \"en\", \"ru\", \"zh\", \"az\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"de",
									"en",
									"ru",
									"zh",
									"az",
								),
							},
						},
						"name": schema.StringAttribute{
							Description: "User's name.",
							Computed:    true,
						},
						"phone": schema.StringAttribute{
							Description: "User's phone.",
							Computed:    true,
						},
						"reseller": schema.Int64Attribute{
							Description: "Services provider ID.",
							Computed:    true,
						},
						"sso_auth": schema.BoolAttribute{
							Description: "SSO authentication flag. If `true` then user can login via SAML SSO.",
							Computed:    true,
						},
						"two_fa": schema.BoolAttribute{
							Description: "Two-step verification:\n- `true` – user enabled two-step verification;\n- `false` – user disabled two-step verification.",
							Computed:    true,
						},
						"user_type": schema.StringAttribute{
							Description: "User's type.\nAvailable values: \"common\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("common"),
							},
						},
					},
				},
			},
		},
	}
}

func (d *IamUsersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *IamUsersDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
