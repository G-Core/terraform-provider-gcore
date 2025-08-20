// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package iam_api_token

import (
	"context"

	"github.com/G-Core/gcore-go/iam"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type IamAPITokenDataSourceModel struct {
	ClientID    types.Int64                                                    `tfsdk:"client_id" path:"clientId,required"`
	TokenID     types.Int64                                                    `tfsdk:"token_id" path:"tokenId,required"`
	Created     types.String                                                   `tfsdk:"created" json:"created,computed"`
	Deleted     types.Bool                                                     `tfsdk:"deleted" json:"deleted,computed"`
	Description types.String                                                   `tfsdk:"description" json:"description,computed"`
	ExpDate     types.String                                                   `tfsdk:"exp_date" json:"exp_date,computed"`
	Expired     types.Bool                                                     `tfsdk:"expired" json:"expired,computed"`
	ID          types.Int64                                                    `tfsdk:"id" json:"id,computed"`
	LastUsage   types.String                                                   `tfsdk:"last_usage" json:"last_usage,computed"`
	Name        types.String                                                   `tfsdk:"name" json:"name,computed"`
	ClientUser  customfield.NestedObject[IamAPITokenClientUserDataSourceModel] `tfsdk:"client_user" json:"client_user,computed"`
}

func (m *IamAPITokenDataSourceModel) toReadParams(_ context.Context) (params iam.APITokenGetParams, diags diag.Diagnostics) {
	params = iam.APITokenGetParams{
		ClientID: m.ClientID.ValueInt64(),
	}

	return
}

type IamAPITokenClientUserDataSourceModel struct {
	ClientID  types.Int64                                                        `tfsdk:"client_id" json:"client_id,computed"`
	Deleted   types.Bool                                                         `tfsdk:"deleted" json:"deleted,computed"`
	Role      customfield.NestedObject[IamAPITokenClientUserRoleDataSourceModel] `tfsdk:"role" json:"role,computed"`
	UserEmail types.String                                                       `tfsdk:"user_email" json:"user_email,computed"`
	UserID    types.Int64                                                        `tfsdk:"user_id" json:"user_id,computed"`
	UserName  types.String                                                       `tfsdk:"user_name" json:"user_name,computed"`
}

type IamAPITokenClientUserRoleDataSourceModel struct {
	ID   types.Int64  `tfsdk:"id" json:"id,computed"`
	Name types.String `tfsdk:"name" json:"name,computed"`
}
