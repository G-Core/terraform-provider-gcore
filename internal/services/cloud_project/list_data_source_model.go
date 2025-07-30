// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_project

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudProjectsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudProjectsItemsDataSourceModel] `json:"results,computed"`
}

type CloudProjectsDataSourceModel struct {
	ClientID       types.Int64                                                     `tfsdk:"client_id" query:"client_id,optional"`
	Name           types.String                                                    `tfsdk:"name" query:"name,optional"`
	IncludeDeleted types.Bool                                                      `tfsdk:"include_deleted" query:"include_deleted,computed_optional"`
	Limit          types.Int64                                                     `tfsdk:"limit" query:"limit,computed_optional"`
	OrderBy        types.String                                                    `tfsdk:"order_by" query:"order_by,computed_optional"`
	MaxItems       types.Int64                                                     `tfsdk:"max_items"`
	Items          customfield.NestedObjectList[CloudProjectsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudProjectsDataSourceModel) toListParams(_ context.Context) (params cloud.ProjectListParams, diags diag.Diagnostics) {
	params = cloud.ProjectListParams{}

	if !m.ClientID.IsNull() {
		params.ClientID = param.NewOpt(m.ClientID.ValueInt64())
	}
	if !m.IncludeDeleted.IsNull() {
		params.IncludeDeleted = param.NewOpt(m.IncludeDeleted.ValueBool())
	}
	if !m.Limit.IsNull() {
		params.Limit = param.NewOpt(m.Limit.ValueInt64())
	}
	if !m.Name.IsNull() {
		params.Name = param.NewOpt(m.Name.ValueString())
	}
	if !m.OrderBy.IsNull() {
		params.OrderBy = cloud.ProjectListParamsOrderBy(m.OrderBy.ValueString())
	}

	return
}

type CloudProjectsItemsDataSourceModel struct {
	ID          types.Int64       `tfsdk:"id" json:"id,computed"`
	ClientID    types.Int64       `tfsdk:"client_id" json:"client_id,computed"`
	CreatedAt   timetypes.RFC3339 `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	IsDefault   types.Bool        `tfsdk:"is_default" json:"is_default,computed"`
	Name        types.String      `tfsdk:"name" json:"name,computed"`
	State       types.String      `tfsdk:"state" json:"state,computed"`
	DeletedAt   timetypes.RFC3339 `tfsdk:"deleted_at" json:"deleted_at,computed" format:"date-time"`
	Description types.String      `tfsdk:"description" json:"description,computed"`
	TaskID      types.String      `tfsdk:"task_id" json:"task_id,computed"`
}
