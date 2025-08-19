// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_audit_log

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudAuditLogsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudAuditLogsItemsDataSourceModel] `json:"results,computed"`
}

type CloudAuditLogsDataSourceModel struct {
	FromTimestamp timetypes.RFC3339                                                `tfsdk:"from_timestamp" query:"from_timestamp,optional" format:"date-time"`
	SearchField   types.String                                                     `tfsdk:"search_field" query:"search_field,optional"`
	ToTimestamp   timetypes.RFC3339                                                `tfsdk:"to_timestamp" query:"to_timestamp,optional" format:"date-time"`
	ActionType    *[]types.String                                                  `tfsdk:"action_type" query:"action_type,optional"`
	APIGroup      *[]types.String                                                  `tfsdk:"api_group" query:"api_group,optional"`
	ProjectID     *[]types.Int64                                                   `tfsdk:"project_id" query:"project_id,optional"`
	RegionID      *[]types.Int64                                                   `tfsdk:"region_id" query:"region_id,optional"`
	ResourceID    *[]types.String                                                  `tfsdk:"resource_id" query:"resource_id,optional"`
	Limit         types.Int64                                                      `tfsdk:"limit" query:"limit,computed_optional"`
	OrderBy       types.String                                                     `tfsdk:"order_by" query:"order_by,computed_optional"`
	Sorting       types.String                                                     `tfsdk:"sorting" query:"sorting,computed_optional"`
	MaxItems      types.Int64                                                      `tfsdk:"max_items"`
	Items         customfield.NestedObjectList[CloudAuditLogsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudAuditLogsDataSourceModel) toListParams(_ context.Context) (params cloud.AuditLogListParams, diags diag.Diagnostics) {
	mActionType := []string{}
	if m.ActionType != nil {
		for _, item := range *m.ActionType {
			mActionType = append(mActionType, string(item.ValueString()))
		}
	}
	mAPIGroup := []string{}
	if m.APIGroup != nil {
		for _, item := range *m.APIGroup {
			mAPIGroup = append(mAPIGroup, string(item.ValueString()))
		}
	}
	mProjectID := []int64{}
	if m.ProjectID != nil {
		for _, item := range *m.ProjectID {
			mProjectID = append(mProjectID, item.ValueInt64())
		}
	}
	mRegionID := []int64{}
	if m.RegionID != nil {
		for _, item := range *m.RegionID {
			mRegionID = append(mRegionID, item.ValueInt64())
		}
	}
	mResourceID := []string{}
	if m.ResourceID != nil {
		for _, item := range *m.ResourceID {
			mResourceID = append(mResourceID, item.ValueString())
		}
	}
	mFromTimestamp, errs := m.FromTimestamp.ValueRFC3339Time()
	diags.Append(errs...)
	mToTimestamp, errs := m.ToTimestamp.ValueRFC3339Time()
	diags.Append(errs...)

	params = cloud.AuditLogListParams{
		ActionType: mActionType,
		APIGroup:   mAPIGroup,
		ProjectID:  mProjectID,
		RegionID:   mRegionID,
		ResourceID: mResourceID,
	}

	if !m.FromTimestamp.IsNull() {
		params.FromTimestamp = param.NewOpt(mFromTimestamp)
	}
	if !m.Limit.IsNull() {
		params.Limit = param.NewOpt(m.Limit.ValueInt64())
	}
	if !m.OrderBy.IsNull() {
		params.OrderBy = cloud.AuditLogListParamsOrderBy(m.OrderBy.ValueString())
	}
	if !m.SearchField.IsNull() {
		params.SearchField = param.NewOpt(m.SearchField.ValueString())
	}
	if !m.Sorting.IsNull() {
		params.Sorting = cloud.AuditLogListParamsSorting(m.Sorting.ValueString())
	}
	if !m.ToTimestamp.IsNull() {
		params.ToTimestamp = param.NewOpt(mToTimestamp)
	}

	return
}

type CloudAuditLogsItemsDataSourceModel struct {
	ID             types.String                                                         `tfsdk:"id" json:"id,computed"`
	Acknowledged   types.Bool                                                           `tfsdk:"acknowledged" json:"acknowledged,computed"`
	ActionData     jsontypes.Normalized                                                 `tfsdk:"action_data" json:"action_data,computed"`
	ActionType     types.String                                                         `tfsdk:"action_type" json:"action_type,computed"`
	APIGroup       types.String                                                         `tfsdk:"api_group" json:"api_group,computed"`
	ClientID       types.Int64                                                          `tfsdk:"client_id" json:"client_id,computed"`
	Email          types.String                                                         `tfsdk:"email" json:"email,computed"`
	IsComplete     types.Bool                                                           `tfsdk:"is_complete" json:"is_complete,computed"`
	IssuedByUserID types.Int64                                                          `tfsdk:"issued_by_user_id" json:"issued_by_user_id,computed"`
	ProjectID      types.Int64                                                          `tfsdk:"project_id" json:"project_id,computed"`
	RegionID       types.Int64                                                          `tfsdk:"region_id" json:"region_id,computed"`
	Resources      customfield.NestedObjectList[CloudAuditLogsResourcesDataSourceModel] `tfsdk:"resources" json:"resources,computed"`
	TaskID         types.String                                                         `tfsdk:"task_id" json:"task_id,computed"`
	Timestamp      timetypes.RFC3339                                                    `tfsdk:"timestamp" json:"timestamp,computed" format:"date-time"`
	TokenID        types.Int64                                                          `tfsdk:"token_id" json:"token_id,computed"`
	TotalPrice     customfield.NestedObject[CloudAuditLogsTotalPriceDataSourceModel]    `tfsdk:"total_price" json:"total_price,computed"`
	UserID         types.Int64                                                          `tfsdk:"user_id" json:"user_id,computed"`
}

type CloudAuditLogsResourcesDataSourceModel struct {
	ResourceID   types.String         `tfsdk:"resource_id" json:"resource_id,computed"`
	ResourceType types.String         `tfsdk:"resource_type" json:"resource_type,computed"`
	ResourceBody jsontypes.Normalized `tfsdk:"resource_body" json:"resource_body,computed"`
	SearchField  types.String         `tfsdk:"search_field" json:"search_field,computed"`
}

type CloudAuditLogsTotalPriceDataSourceModel struct {
	CurrencyCode  types.String  `tfsdk:"currency_code" json:"currency_code,computed"`
	PricePerHour  types.Float64 `tfsdk:"price_per_hour" json:"price_per_hour,computed"`
	PricePerMonth types.Float64 `tfsdk:"price_per_month" json:"price_per_month,computed"`
	PriceStatus   types.String  `tfsdk:"price_status" json:"price_status,computed"`
}
