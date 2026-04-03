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
	ID          types.Int64                           `tfsdk:"id" path:"project_id,computed"`
	ProjectID   types.Int64                           `tfsdk:"project_id" path:"project_id,optional"`
	ClientID    types.Int64                           `tfsdk:"client_id" json:"client_id,computed"`
	CreatedAt   timetypes.RFC3339                     `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	DeletedAt   timetypes.RFC3339                     `tfsdk:"deleted_at" json:"deleted_at,computed" format:"date-time"`
	Description types.String                          `tfsdk:"description" json:"description,computed"`
	IsDefault   types.Bool                            `tfsdk:"is_default" json:"is_default,computed"`
	Name        types.String                          `tfsdk:"name" json:"name,computed"`
	State       types.String                          `tfsdk:"state" json:"state,computed"`
	TaskID      types.String                          `tfsdk:"task_id" json:"task_id,computed"`
	FindOneBy   *CloudProjectFindOneByDataSourceModel `tfsdk:"find_one_by"`
}

func (m *CloudProjectDataSourceModel) toReadParams(_ context.Context) (params cloud.ProjectGetParams, diags diag.Diagnostics) {
	params = cloud.ProjectGetParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}

	return
}

func (m *CloudProjectDataSourceModel) toListParams(_ context.Context) (params cloud.ProjectListParams, diags diag.Diagnostics) {
	params = cloud.ProjectListParams{}

	if !m.FindOneBy.IncludeDeleted.IsNull() {
		params.IncludeDeleted = param.NewOpt(m.FindOneBy.IncludeDeleted.ValueBool())
	}
	if !m.FindOneBy.Limit.IsNull() {
		params.Limit = param.NewOpt(m.FindOneBy.Limit.ValueInt64())
	}
	if !m.FindOneBy.Name.IsNull() {
		params.Name = param.NewOpt(m.FindOneBy.Name.ValueString())
	}
	if !m.FindOneBy.OrderBy.IsNull() {
		params.OrderBy = cloud.ProjectListParamsOrderBy(m.FindOneBy.OrderBy.ValueString())
	}

	return
}

type CloudProjectFindOneByDataSourceModel struct {
	IncludeDeleted types.Bool   `tfsdk:"include_deleted" query:"include_deleted,computed_optional"`
	Limit          types.Int64  `tfsdk:"limit" query:"limit,computed_optional"`
	Name           types.String `tfsdk:"name" query:"name,optional"`
	OrderBy        types.String `tfsdk:"order_by" query:"order_by,computed_optional"`
}
