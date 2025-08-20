// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package iam_user

import (
	"context"

	"github.com/G-Core/gcore-go/iam"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type IamUsersResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[IamUsersItemsDataSourceModel] `json:"results,computed"`
}

type IamUsersDataSourceModel struct {
	Limit    types.Int64                                                `tfsdk:"limit" query:"limit,optional"`
	MaxItems types.Int64                                                `tfsdk:"max_items"`
	Items    customfield.NestedObjectList[IamUsersItemsDataSourceModel] `tfsdk:"items"`
}

func (m *IamUsersDataSourceModel) toListParams(_ context.Context) (params iam.UserListParams, diags diag.Diagnostics) {
	params = iam.UserListParams{}

	if !m.Limit.IsNull() {
		params.Limit = param.NewOpt(m.Limit.ValueInt64())
	}

	return
}

type IamUsersItemsDataSourceModel struct {
	ID        types.Int64                                                 `tfsdk:"id" json:"id,computed"`
	Activated types.Bool                                                  `tfsdk:"activated" json:"activated,computed"`
	AuthTypes customfield.List[types.String]                              `tfsdk:"auth_types" json:"auth_types,computed"`
	Client    types.Float64                                               `tfsdk:"client" json:"client,computed"`
	Company   types.String                                                `tfsdk:"company" json:"company,computed"`
	Deleted   types.Bool                                                  `tfsdk:"deleted" json:"deleted,computed"`
	Email     types.String                                                `tfsdk:"email" json:"email,computed"`
	Groups    customfield.NestedObjectList[IamUsersGroupsDataSourceModel] `tfsdk:"groups" json:"groups,computed"`
	Lang      types.String                                                `tfsdk:"lang" json:"lang,computed"`
	Name      types.String                                                `tfsdk:"name" json:"name,computed"`
	Phone     types.String                                                `tfsdk:"phone" json:"phone,computed"`
	Reseller  types.Int64                                                 `tfsdk:"reseller" json:"reseller,computed"`
	SSOAuth   types.Bool                                                  `tfsdk:"sso_auth" json:"sso_auth,computed"`
	TwoFa     types.Bool                                                  `tfsdk:"two_fa" json:"two_fa,computed"`
	UserType  types.String                                                `tfsdk:"user_type" json:"user_type,computed"`
}

type IamUsersGroupsDataSourceModel struct {
	ID   types.Int64  `tfsdk:"id" json:"id,computed"`
	Name types.String `tfsdk:"name" json:"name,computed"`
}
