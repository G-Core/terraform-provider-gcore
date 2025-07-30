// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package iam_api_token

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
)

type IamAPITokenModel struct {
	ID          types.String                `tfsdk:"id" json:"-,computed"`
	Token       types.String                `tfsdk:"token" json:"token,computed"`
	ClientID    types.Int64                 `tfsdk:"client_id" path:"clientId,required"`
	TokenID     types.Int64                 `tfsdk:"token_id" path:"tokenId,optional"`
	ExpDate     types.String                `tfsdk:"exp_date" json:"exp_date,required"`
	Name        types.String                `tfsdk:"name" json:"name,required"`
	ClientUser  *IamAPITokenClientUserModel `tfsdk:"client_user" json:"client_user,required"`
	Description types.String                `tfsdk:"description" json:"description,optional"`
	Created     types.String                `tfsdk:"created" json:"created,computed"`
	Deleted     types.Bool                  `tfsdk:"deleted" json:"deleted,computed"`
	Expired     types.Bool                  `tfsdk:"expired" json:"expired,computed"`
	LastUsage   types.String                `tfsdk:"last_usage" json:"last_usage,computed"`
}

func (m IamAPITokenModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m IamAPITokenModel) MarshalJSONForUpdate(state IamAPITokenModel) (data []byte, err error) {
	return apijson.MarshalForUpdate(m, state)
}

type IamAPITokenClientUserModel struct {
	Role *IamAPITokenClientUserRoleModel `tfsdk:"role" json:"role,optional"`
}

type IamAPITokenClientUserRoleModel struct {
	ID   types.Int64  `tfsdk:"id" json:"id,optional"`
	Name types.String `tfsdk:"name" json:"name,optional"`
}
