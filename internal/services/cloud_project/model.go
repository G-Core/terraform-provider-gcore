// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_project

import (
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudProjectModel struct {
	ID          types.Int64       `tfsdk:"id" json:"id,computed"`
	Name        types.String      `tfsdk:"name" json:"name,required"`
	Description types.String      `tfsdk:"description" json:"description,computed_optional"`
	ClientID    types.Int64       `tfsdk:"client_id" json:"client_id,computed"`
	CreatedAt   timetypes.RFC3339 `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	IsDefault   types.Bool        `tfsdk:"is_default" json:"is_default,computed"`
	State       types.String      `tfsdk:"state" json:"state,computed"`
}

func (m CloudProjectModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudProjectModel) MarshalJSONForUpdate(state CloudProjectModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}
