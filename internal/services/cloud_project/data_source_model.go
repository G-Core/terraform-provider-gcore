// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_project

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudProjectDataSourceModel struct {
	ProjectID   types.Int64       `tfsdk:"project_id" path:"project_id,required"`
	ClientID    types.Int64       `tfsdk:"client_id" json:"client_id,computed"`
	CreatedAt   timetypes.RFC3339 `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	DeletedAt   timetypes.RFC3339 `tfsdk:"deleted_at" json:"deleted_at,computed" format:"date-time"`
	Description types.String      `tfsdk:"description" json:"description,computed"`
	ID          types.Int64       `tfsdk:"id" json:"id,computed"`
	IsDefault   types.Bool        `tfsdk:"is_default" json:"is_default,computed"`
	Name        types.String      `tfsdk:"name" json:"name,computed"`
	State       types.String      `tfsdk:"state" json:"state,computed"`
	TaskID      types.String      `tfsdk:"task_id" json:"task_id,computed"`
}

func (m *CloudProjectDataSourceModel) toReadParams(_ context.Context) (params cloud.ProjectGetParams, diags diag.Diagnostics) {
	params = cloud.ProjectGetParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}

	return
}
