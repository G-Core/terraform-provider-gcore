// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package iam_user

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type IamUserModel struct {
	ID             types.Int64                                              `tfsdk:"id" json:"id,computed"`
	UserID         types.Int64                                              `tfsdk:"user_id" path:"userId,required"`
	ClientID       types.Int64                                              `tfsdk:"client_id" path:"clientId,optional"`
	Company        types.String                                             `tfsdk:"company" json:"company,optional"`
	Email          types.String                                             `tfsdk:"email" json:"email,optional"`
	Lang           types.String                                             `tfsdk:"lang" json:"lang,optional"`
	Name           types.String                                             `tfsdk:"name" json:"name,optional"`
	Phone          types.String                                             `tfsdk:"phone" json:"phone,optional"`
	AuthTypes      *[]types.String                                          `tfsdk:"auth_types" json:"auth_types,optional"`
	Groups         *[]*IamUserGroupsModel                                   `tfsdk:"groups" json:"groups,optional"`
	Activated      types.Bool                                               `tfsdk:"activated" json:"activated,computed"`
	Client         types.Float64                                            `tfsdk:"client" json:"client,computed"`
	Deleted        types.Bool                                               `tfsdk:"deleted" json:"deleted,computed"`
	IsActive       types.Bool                                               `tfsdk:"is_active" json:"is_active,computed"`
	Reseller       types.Int64                                              `tfsdk:"reseller" json:"reseller,computed"`
	SSOAuth        types.Bool                                               `tfsdk:"sso_auth" json:"sso_auth,computed"`
	TwoFa          types.Bool                                               `tfsdk:"two_fa" json:"two_fa,computed"`
	UserType       types.String                                             `tfsdk:"user_type" json:"user_type,computed"`
	ClientAndRoles customfield.NestedObjectList[IamUserClientAndRolesModel] `tfsdk:"client_and_roles" json:"client_and_roles,computed"`
}

func (m IamUserModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m IamUserModel) MarshalJSONForUpdate(state IamUserModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type IamUserGroupsModel struct {
	ID   types.Int64  `tfsdk:"id" json:"id,optional"`
	Name types.String `tfsdk:"name" json:"name,optional"`
}

type IamUserClientAndRolesModel struct {
	ClientCompanyName types.String                   `tfsdk:"client_company_name" json:"client_company_name,computed"`
	ClientID          types.Int64                    `tfsdk:"client_id" json:"client_id,computed"`
	UserID            types.Int64                    `tfsdk:"user_id" json:"user_id,computed"`
	UserRoles         customfield.List[types.String] `tfsdk:"user_roles" json:"user_roles,computed"`
}
