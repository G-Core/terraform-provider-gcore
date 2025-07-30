// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package iam_user

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type IamUserDataSourceModel struct {
	UserID         types.Int64                                                        `tfsdk:"user_id" path:"userId,required"`
	Activated      types.Bool                                                         `tfsdk:"activated" json:"activated,computed"`
	Client         types.Float64                                                      `tfsdk:"client" json:"client,computed"`
	Company        types.String                                                       `tfsdk:"company" json:"company,computed"`
	Deleted        types.Bool                                                         `tfsdk:"deleted" json:"deleted,computed"`
	Email          types.String                                                       `tfsdk:"email" json:"email,computed"`
	ID             types.Int64                                                        `tfsdk:"id" json:"id,computed"`
	IsActive       types.Bool                                                         `tfsdk:"is_active" json:"is_active,computed"`
	Lang           types.String                                                       `tfsdk:"lang" json:"lang,computed"`
	Name           types.String                                                       `tfsdk:"name" json:"name,computed"`
	Phone          types.String                                                       `tfsdk:"phone" json:"phone,computed"`
	Reseller       types.Int64                                                        `tfsdk:"reseller" json:"reseller,computed"`
	SSOAuth        types.Bool                                                         `tfsdk:"sso_auth" json:"sso_auth,computed"`
	TwoFa          types.Bool                                                         `tfsdk:"two_fa" json:"two_fa,computed"`
	UserType       types.String                                                       `tfsdk:"user_type" json:"user_type,computed"`
	AuthTypes      customfield.List[types.String]                                     `tfsdk:"auth_types" json:"auth_types,computed"`
	ClientAndRoles customfield.NestedObjectList[IamUserClientAndRolesDataSourceModel] `tfsdk:"client_and_roles" json:"client_and_roles,computed"`
	Groups         customfield.NestedObjectList[IamUserGroupsDataSourceModel]         `tfsdk:"groups" json:"groups,computed"`
}

type IamUserClientAndRolesDataSourceModel struct {
	ClientCompanyName types.String                   `tfsdk:"client_company_name" json:"client_company_name,computed"`
	ClientID          types.Int64                    `tfsdk:"client_id" json:"client_id,computed"`
	UserID            types.Int64                    `tfsdk:"user_id" json:"user_id,computed"`
	UserRoles         customfield.List[types.String] `tfsdk:"user_roles" json:"user_roles,computed"`
}

type IamUserGroupsDataSourceModel struct {
	ID   types.Int64  `tfsdk:"id" json:"id,computed"`
	Name types.String `tfsdk:"name" json:"name,computed"`
}
