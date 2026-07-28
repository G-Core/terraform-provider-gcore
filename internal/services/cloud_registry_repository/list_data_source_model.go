// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_registry_repository

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudRegistryRepositoriesResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudRegistryRepositoriesItemsDataSourceModel] `json:"results,computed"`
}

type CloudRegistryRepositoriesDataSourceModel struct {
	RegistryID types.Int64                                                                 `tfsdk:"registry_id" path:"registry_id,required"`
	ProjectID  types.Int64                                                                 `tfsdk:"project_id" path:"project_id,optional"`
	RegionID   types.Int64                                                                 `tfsdk:"region_id" path:"region_id,optional"`
	MaxItems   types.Int64                                                                 `tfsdk:"max_items"`
	Items      customfield.NestedObjectList[CloudRegistryRepositoriesItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudRegistryRepositoriesDataSourceModel) toListParams(_ context.Context) (params cloud.RegistryRepositoryListParams, diags diag.Diagnostics) {
	params = cloud.RegistryRepositoryListParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}

	return
}

type CloudRegistryRepositoriesItemsDataSourceModel struct {
	ID            types.Int64       `tfsdk:"id" json:"id,computed"`
	ArtifactCount types.Int64       `tfsdk:"artifact_count" json:"artifact_count,computed"`
	CreatedAt     timetypes.RFC3339 `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	Name          types.String      `tfsdk:"name" json:"name,computed"`
	PullCount     types.Int64       `tfsdk:"pull_count" json:"pull_count,computed"`
	RegistryID    types.Int64       `tfsdk:"registry_id" json:"registry_id,computed"`
	UpdatedAt     timetypes.RFC3339 `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
}
