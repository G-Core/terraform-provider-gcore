// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_file_share_access_rule

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudFileShareAccessRulesResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudFileShareAccessRulesItemsDataSourceModel] `json:"results,computed"`
}

type CloudFileShareAccessRulesDataSourceModel struct {
	FileShareID types.String                                                                `tfsdk:"file_share_id" path:"file_share_id,required"`
	ProjectID   types.Int64                                                                 `tfsdk:"project_id" path:"project_id,optional"`
	RegionID    types.Int64                                                                 `tfsdk:"region_id" path:"region_id,optional"`
	MaxItems    types.Int64                                                                 `tfsdk:"max_items"`
	Items       customfield.NestedObjectList[CloudFileShareAccessRulesItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudFileShareAccessRulesDataSourceModel) toListParams(_ context.Context) (params cloud.FileShareAccessRuleListParams, diags diag.Diagnostics) {
	params = cloud.FileShareAccessRuleListParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}

	return
}

type CloudFileShareAccessRulesItemsDataSourceModel struct {
	ID          types.String `tfsdk:"id" json:"id,computed"`
	AccessLevel types.String `tfsdk:"access_level" json:"access_level,computed"`
	AccessTo    types.String `tfsdk:"access_to" json:"access_to,computed"`
	State       types.String `tfsdk:"state" json:"state,computed"`
}
