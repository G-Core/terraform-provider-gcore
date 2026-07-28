// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_registry

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudRegistriesResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudRegistriesItemsDataSourceModel] `json:"results,computed"`
}

type CloudRegistriesDataSourceModel struct {
	ProjectID types.Int64                                                       `tfsdk:"project_id" path:"project_id,optional"`
	RegionID  types.Int64                                                       `tfsdk:"region_id" path:"region_id,optional"`
	MaxItems  types.Int64                                                       `tfsdk:"max_items"`
	Items     customfield.NestedObjectList[CloudRegistriesItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudRegistriesDataSourceModel) toListParams(_ context.Context) (params cloud.RegistryListParams, diags diag.Diagnostics) {
	params = cloud.RegistryListParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}

	return
}

type CloudRegistriesItemsDataSourceModel struct {
	ID           types.Int64       `tfsdk:"id" json:"id,computed"`
	CreatedAt    timetypes.RFC3339 `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	Name         types.String      `tfsdk:"name" json:"name,computed"`
	RepoCount    types.Int64       `tfsdk:"repo_count" json:"repo_count,computed"`
	StorageLimit types.Int64       `tfsdk:"storage_limit" json:"storage_limit,computed"`
	StorageUsed  types.Int64       `tfsdk:"storage_used" json:"storage_used,computed"`
	UpdatedAt    timetypes.RFC3339 `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	URL          types.String      `tfsdk:"url" json:"url,computed"`
}
