// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_registry_user

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudRegistryUsersResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudRegistryUsersItemsDataSourceModel] `json:"results,computed"`
}

type CloudRegistryUsersDataSourceModel struct {
	RegistryID types.Int64                                                          `tfsdk:"registry_id" path:"registry_id,required"`
	ProjectID  types.Int64                                                          `tfsdk:"project_id" path:"project_id,optional"`
	RegionID   types.Int64                                                          `tfsdk:"region_id" path:"region_id,optional"`
	MaxItems   types.Int64                                                          `tfsdk:"max_items"`
	Items      customfield.NestedObjectList[CloudRegistryUsersItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudRegistryUsersDataSourceModel) toListParams(_ context.Context) (params cloud.RegistryUserListParams, diags diag.Diagnostics) {
	params = cloud.RegistryUserListParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}

	return
}

type CloudRegistryUsersItemsDataSourceModel struct {
	ID        types.Int64       `tfsdk:"id" json:"id,computed"`
	CreatedAt timetypes.RFC3339 `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	Duration  types.Int64       `tfsdk:"duration" json:"duration,computed"`
	ExpiresAt timetypes.RFC3339 `tfsdk:"expires_at" json:"expires_at,computed" format:"date-time"`
	Name      types.String      `tfsdk:"name" json:"name,computed"`
	ReadOnly  types.Bool        `tfsdk:"read_only" json:"read_only,computed"`
}
