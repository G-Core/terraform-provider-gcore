// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_registry

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudRegistryDataSourceModel struct {
	ProjectID    types.Int64       `tfsdk:"project_id" path:"project_id,required"`
	RegionID     types.Int64       `tfsdk:"region_id" path:"region_id,required"`
	RegistryID   types.Int64       `tfsdk:"registry_id" path:"registry_id,required"`
	CreatedAt    timetypes.RFC3339 `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	ID           types.Int64       `tfsdk:"id" json:"id,computed"`
	Name         types.String      `tfsdk:"name" json:"name,computed"`
	RepoCount    types.Int64       `tfsdk:"repo_count" json:"repo_count,computed"`
	StorageLimit types.Int64       `tfsdk:"storage_limit" json:"storage_limit,computed"`
	StorageUsed  types.Int64       `tfsdk:"storage_used" json:"storage_used,computed"`
	UpdatedAt    timetypes.RFC3339 `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	URL          types.String      `tfsdk:"url" json:"url,computed"`
}

func (m *CloudRegistryDataSourceModel) toReadParams(_ context.Context) (params cloud.RegistryGetParams, diags diag.Diagnostics) {
	params = cloud.RegistryGetParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}

	return
}
