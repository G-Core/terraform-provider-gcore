// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_user_role_assignment

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudUserRoleAssignmentsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudUserRoleAssignmentsItemsDataSourceModel] `json:"results,computed"`
}

type CloudUserRoleAssignmentsDataSourceModel struct {
	ProjectID types.Int64                                                                `tfsdk:"project_id" query:"project_id,optional"`
	UserID    types.Int64                                                                `tfsdk:"user_id" query:"user_id,optional"`
	Limit     types.Int64                                                                `tfsdk:"limit" query:"limit,computed_optional"`
	MaxItems  types.Int64                                                                `tfsdk:"max_items"`
	Items     customfield.NestedObjectList[CloudUserRoleAssignmentsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudUserRoleAssignmentsDataSourceModel) toListParams(_ context.Context) (params cloud.UserRoleAssignmentListParams, diags diag.Diagnostics) {
	params = cloud.UserRoleAssignmentListParams{}

	if !m.Limit.IsNull() {
		params.Limit = param.NewOpt(m.Limit.ValueInt64())
	}
	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.UserID.IsNull() {
		params.UserID = param.NewOpt(m.UserID.ValueInt64())
	}

	return
}

type CloudUserRoleAssignmentsItemsDataSourceModel struct {
	ID         types.Int64       `tfsdk:"id" json:"id,computed"`
	AssignedBy types.Int64       `tfsdk:"assigned_by" json:"assigned_by,computed"`
	ClientID   types.Int64       `tfsdk:"client_id" json:"client_id,computed"`
	CreatedAt  timetypes.RFC3339 `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	ProjectID  types.Int64       `tfsdk:"project_id" json:"project_id,computed"`
	Role       types.String      `tfsdk:"role" json:"role,computed"`
	UpdatedAt  timetypes.RFC3339 `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	UserID     types.Int64       `tfsdk:"user_id" json:"user_id,computed"`
}
