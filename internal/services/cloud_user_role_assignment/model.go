// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_user_role_assignment

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
)

type CloudUserRoleAssignmentModel struct {
	ID           types.Int64       `tfsdk:"id" json:"id,computed"`
	Role         types.String      `tfsdk:"role" json:"role,required"`
	UserID       types.Int64       `tfsdk:"user_id" json:"user_id,required"`
	ClientID     types.Int64       `tfsdk:"client_id" json:"client_id,optional"`
	ProjectID    types.Int64       `tfsdk:"project_id" json:"project_id,optional"`
	AssignedBy   types.Int64       `tfsdk:"assigned_by" json:"assigned_by,computed"`
	AssignmentID types.Int64       `tfsdk:"assignment_id" json:"assignment_id,computed"`
	CreatedAt    timetypes.RFC3339 `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	UpdatedAt    timetypes.RFC3339 `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
}

func (m CloudUserRoleAssignmentModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudUserRoleAssignmentModel) MarshalJSONForUpdate(state CloudUserRoleAssignmentModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}
