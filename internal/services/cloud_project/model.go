// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_project

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
)

type CloudProjectModel struct {
	ID          types.Int64       `tfsdk:"id" json:"id,computed"`
	Name        types.String      `tfsdk:"name" json:"name,required"`
	ClientID    types.Int64       `tfsdk:"client_id" json:"client_id,optional"`
	Description types.String      `tfsdk:"description" json:"description,optional"`
	State       types.String      `tfsdk:"state" json:"state,computed_optional"`
	CreatedAt   timetypes.RFC3339 `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	DeletedAt   timetypes.RFC3339 `tfsdk:"deleted_at" json:"deleted_at,computed" format:"date-time"`
	IsDefault   types.Bool        `tfsdk:"is_default" json:"is_default,computed"`
	TaskID      types.String      `tfsdk:"task_id" json:"task_id,computed"`
}

func (m CloudProjectModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudProjectModel) MarshalJSONForUpdate(state CloudProjectModel) (data []byte, err error) {
	return apijson.MarshalForUpdate(m, state)
}
