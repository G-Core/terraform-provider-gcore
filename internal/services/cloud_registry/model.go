// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_registry

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
)

type CloudRegistryModel struct {
	ID           types.Int64       `tfsdk:"id" json:"id,computed"`
	ProjectID    types.Int64       `tfsdk:"project_id" path:"project_id,optional"`
	RegionID     types.Int64       `tfsdk:"region_id" path:"region_id,optional"`
	Name         types.String      `tfsdk:"name" json:"name,required"`
	StorageLimit types.Int64       `tfsdk:"storage_limit" json:"storage_limit,computed_optional"`
	CreatedAt    timetypes.RFC3339 `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	RepoCount    types.Int64       `tfsdk:"repo_count" json:"repo_count,computed"`
	StorageUsed  types.Int64       `tfsdk:"storage_used" json:"storage_used,computed"`
	UpdatedAt    timetypes.RFC3339 `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	URL          types.String      `tfsdk:"url" json:"url,computed"`
}

func (m CloudRegistryModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudRegistryModel) MarshalJSONForUpdate(state CloudRegistryModel) (data []byte, err error) {
	return apijson.MarshalForUpdate(m, state)
}
