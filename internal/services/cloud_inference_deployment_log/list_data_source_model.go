// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_deployment_log

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudInferenceDeploymentLogsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudInferenceDeploymentLogsItemsDataSourceModel] `json:"results,computed"`
}

type CloudInferenceDeploymentLogsDataSourceModel struct {
	DeploymentName types.String                                                                   `tfsdk:"deployment_name" path:"deployment_name,required"`
	ProjectID      types.Int64                                                                    `tfsdk:"project_id" path:"project_id,optional"`
	RegionID       types.Int64                                                                    `tfsdk:"region_id" query:"region_id,optional"`
	Limit          types.Int64                                                                    `tfsdk:"limit" query:"limit,computed_optional"`
	OrderBy        types.String                                                                   `tfsdk:"order_by" query:"order_by,computed_optional"`
	MaxItems       types.Int64                                                                    `tfsdk:"max_items"`
	Items          customfield.NestedObjectList[CloudInferenceDeploymentLogsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudInferenceDeploymentLogsDataSourceModel) toListParams(_ context.Context) (params cloud.InferenceDeploymentLogListParams, diags diag.Diagnostics) {
	params = cloud.InferenceDeploymentLogListParams{}

	if !m.Limit.IsNull() {
		params.Limit = param.NewOpt(m.Limit.ValueInt64())
	}
	if !m.OrderBy.IsNull() {
		params.OrderBy = cloud.InferenceDeploymentLogListParamsOrderBy(m.OrderBy.ValueString())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}
	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}

	return
}

type CloudInferenceDeploymentLogsItemsDataSourceModel struct {
	Message  types.String      `tfsdk:"message" json:"message,computed"`
	Pod      types.String      `tfsdk:"pod" json:"pod,computed"`
	RegionID types.Int64       `tfsdk:"region_id" json:"region_id,computed"`
	Time     timetypes.RFC3339 `tfsdk:"time" json:"time,computed" format:"date-time"`
}
