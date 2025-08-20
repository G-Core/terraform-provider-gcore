// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package iam_user

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ resource.ResourceWithConfigValidators = (*IamUserResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:   "User's ID.",
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"user_id": schema.Int64Attribute{
				Required:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"client_id": schema.Int64Attribute{
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"company": schema.StringAttribute{
				Description: "User's company.",
				Optional:    true,
			},
			"email": schema.StringAttribute{
				Description: "User's email address.",
				Optional:    true,
			},
			"lang": schema.StringAttribute{
				Description: "User's language.\nDefines language of the control panel and email messages.\nAvailable values: \"de\", \"en\", \"ru\", \"zh\", \"az\".",
				Optional:    true,
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
				Optional:    true,
			},
			"phone": schema.StringAttribute{
				Description: "User's phone.",
				Optional:    true,
			},
			"auth_types": schema.ListAttribute{
				Description: "System field. List of auth types available for the account.",
				Optional:    true,
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
				ElementType: types.StringType,
			},
			"groups": schema.ListNestedAttribute{
				Description: "User's group in the current account.\nIAM supports 5 groups:\n- Users\n- Administrators\n- Engineers\n- Purge and Prefetch only (API)\n- Purge and Prefetch only (API+Web)",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
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
			"activated": schema.BoolAttribute{
				Description: "Email confirmation:\n- `true` – user confirmed the email;\n- `false` – user did not confirm the email.",
				Computed:    true,
			},
			"client": schema.Float64Attribute{
				Description: "User's account ID.",
				Computed:    true,
			},
			"deleted": schema.BoolAttribute{
				Description: "Deletion flag. If `true` then user was deleted.",
				Computed:    true,
			},
			"is_active": schema.BoolAttribute{
				Description: "User activity flag.",
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
				Description: "User's type.\nAvailable values: \"common\", \"reseller\", \"seller\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"common",
						"reseller",
						"seller",
					),
				},
			},
			"client_and_roles": schema.ListNestedAttribute{
				Description: "List of user's clients. User can access to one or more clients.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[IamUserClientAndRolesModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"client_company_name": schema.StringAttribute{
							Computed: true,
						},
						"client_id": schema.Int64Attribute{
							Computed: true,
						},
						"user_id": schema.Int64Attribute{
							Description: "User's ID.",
							Computed:    true,
						},
						"user_roles": schema.ListAttribute{
							Description: "User role in this client.",
							Computed:    true,
							CustomType:  customfield.NewListType[types.String](ctx),
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func (r *IamUserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *IamUserResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
