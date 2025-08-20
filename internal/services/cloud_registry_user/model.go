// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_registry_user

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
)

type CloudRegistryUserModel struct {
	ID         types.Int64       `tfsdk:"id" json:"id,computed"`
	RegistryID types.Int64       `tfsdk:"registry_id" path:"registry_id,required"`
	ProjectID  types.Int64       `tfsdk:"project_id" path:"project_id,optional"`
	RegionID   types.Int64       `tfsdk:"region_id" path:"region_id,optional"`
	Name       types.String      `tfsdk:"name" json:"name,required"`
	Secret     types.String      `tfsdk:"secret" json:"secret,computed_optional"`
	Duration   types.Int64       `tfsdk:"duration" json:"duration,required"`
	ReadOnly   types.Bool        `tfsdk:"read_only" json:"read_only,computed_optional"`
	CreatedAt  timetypes.RFC3339 `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	ExpiresAt  timetypes.RFC3339 `tfsdk:"expires_at" json:"expires_at,computed" format:"date-time"`
}

func (m CloudRegistryUserModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudRegistryUserModel) MarshalJSONForUpdate(state CloudRegistryUserModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}
