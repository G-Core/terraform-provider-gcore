// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_file_share_access_rule

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
)

type CloudFileShareAccessRuleModel struct {
	ID          types.String `tfsdk:"id" json:"id,computed"`
	FileShareID types.String `tfsdk:"file_share_id" path:"file_share_id,required"`
	ProjectID   types.Int64  `tfsdk:"project_id" path:"project_id,optional"`
	RegionID    types.Int64  `tfsdk:"region_id" path:"region_id,optional"`
	AccessMode  types.String `tfsdk:"access_mode" json:"access_mode,required"`
	IPAddress   types.String `tfsdk:"ip_address" json:"ip_address,required"`
	AccessLevel types.String `tfsdk:"access_level" json:"access_level,computed"`
	AccessTo    types.String `tfsdk:"access_to" json:"access_to,computed"`
	State       types.String `tfsdk:"state" json:"state,computed"`
}

func (m CloudFileShareAccessRuleModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudFileShareAccessRuleModel) MarshalJSONForUpdate(state CloudFileShareAccessRuleModel) (data []byte, err error) {
	return apijson.MarshalForUpdate(m, state)
}
