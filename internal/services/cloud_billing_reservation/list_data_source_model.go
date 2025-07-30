// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_billing_reservation

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudBillingReservationsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudBillingReservationsItemsDataSourceModel] `json:"results,computed"`
}

type CloudBillingReservationsDataSourceModel struct {
	ActivatedFrom   timetypes.RFC3339                                                          `tfsdk:"activated_from" query:"activated_from,optional" format:"date"`
	ActivatedTo     timetypes.RFC3339                                                          `tfsdk:"activated_to" query:"activated_to,optional" format:"date"`
	CreatedFrom     timetypes.RFC3339                                                          `tfsdk:"created_from" query:"created_from,optional" format:"date-time"`
	CreatedTo       timetypes.RFC3339                                                          `tfsdk:"created_to" query:"created_to,optional" format:"date-time"`
	DeactivatedFrom timetypes.RFC3339                                                          `tfsdk:"deactivated_from" query:"deactivated_from,optional" format:"date"`
	DeactivatedTo   timetypes.RFC3339                                                          `tfsdk:"deactivated_to" query:"deactivated_to,optional" format:"date"`
	MetricName      types.String                                                               `tfsdk:"metric_name" query:"metric_name,optional"`
	RegionID        types.Int64                                                                `tfsdk:"region_id" query:"region_id,optional"`
	Status          *[]types.String                                                            `tfsdk:"status" query:"status,optional"`
	Limit           types.Int64                                                                `tfsdk:"limit" query:"limit,computed_optional"`
	OrderBy         types.String                                                               `tfsdk:"order_by" query:"order_by,computed_optional"`
	MaxItems        types.Int64                                                                `tfsdk:"max_items"`
	Items           customfield.NestedObjectList[CloudBillingReservationsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudBillingReservationsDataSourceModel) toListParams(_ context.Context) (params cloud.BillingReservationListParams, diags diag.Diagnostics) {
	mStatus := []string{}
	for _, item := range *m.Status {
		mStatus = append(mStatus, string(item.ValueString()))
	}
	mActivatedFrom, errs := m.ActivatedFrom.ValueRFC3339Time()
	diags.Append(errs...)
	mActivatedTo, errs := m.ActivatedTo.ValueRFC3339Time()
	diags.Append(errs...)
	mCreatedFrom, errs := m.CreatedFrom.ValueRFC3339Time()
	diags.Append(errs...)
	mCreatedTo, errs := m.CreatedTo.ValueRFC3339Time()
	diags.Append(errs...)
	mDeactivatedFrom, errs := m.DeactivatedFrom.ValueRFC3339Time()
	diags.Append(errs...)
	mDeactivatedTo, errs := m.DeactivatedTo.ValueRFC3339Time()
	diags.Append(errs...)

	params = cloud.BillingReservationListParams{
		Status: mStatus,
	}

	if !m.ActivatedFrom.IsNull() {
		params.ActivatedFrom = param.NewOpt(mActivatedFrom)
	}
	if !m.ActivatedTo.IsNull() {
		params.ActivatedTo = param.NewOpt(mActivatedTo)
	}
	if !m.CreatedFrom.IsNull() {
		params.CreatedFrom = param.NewOpt(mCreatedFrom)
	}
	if !m.CreatedTo.IsNull() {
		params.CreatedTo = param.NewOpt(mCreatedTo)
	}
	if !m.DeactivatedFrom.IsNull() {
		params.DeactivatedFrom = param.NewOpt(mDeactivatedFrom)
	}
	if !m.DeactivatedTo.IsNull() {
		params.DeactivatedTo = param.NewOpt(mDeactivatedTo)
	}
	if !m.Limit.IsNull() {
		params.Limit = param.NewOpt(m.Limit.ValueInt64())
	}
	if !m.MetricName.IsNull() {
		params.MetricName = param.NewOpt(m.MetricName.ValueString())
	}
	if !m.OrderBy.IsNull() {
		params.OrderBy = cloud.BillingReservationListParamsOrderBy(m.OrderBy.ValueString())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}

	return
}

type CloudBillingReservationsItemsDataSourceModel struct {
	ID                         types.Int64                                                                    `tfsdk:"id" json:"id,computed"`
	ActiveFrom                 timetypes.RFC3339                                                              `tfsdk:"active_from" json:"active_from,computed" format:"date"`
	ActiveTo                   timetypes.RFC3339                                                              `tfsdk:"active_to" json:"active_to,computed" format:"date"`
	ActivityPeriod             types.String                                                                   `tfsdk:"activity_period" json:"activity_period,computed"`
	ActivityPeriodLength       types.Int64                                                                    `tfsdk:"activity_period_length" json:"activity_period_length,computed"`
	AmountPrices               customfield.NestedObject[CloudBillingReservationsAmountPricesDataSourceModel]  `tfsdk:"amount_prices" json:"amount_prices,computed"`
	BillingPlanID              types.Int64                                                                    `tfsdk:"billing_plan_id" json:"billing_plan_id,computed"`
	CreatedAt                  timetypes.RFC3339                                                              `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	Error                      types.String                                                                   `tfsdk:"error" json:"error,computed"`
	Eta                        timetypes.RFC3339                                                              `tfsdk:"eta" json:"eta,computed" format:"date"`
	IsExpirationMessageVisible types.Bool                                                                     `tfsdk:"is_expiration_message_visible" json:"is_expiration_message_visible,computed"`
	Name                       types.String                                                                   `tfsdk:"name" json:"name,computed"`
	NextStatuses               customfield.List[types.String]                                                 `tfsdk:"next_statuses" json:"next_statuses,computed"`
	RegionID                   types.Int64                                                                    `tfsdk:"region_id" json:"region_id,computed"`
	RegionName                 types.String                                                                   `tfsdk:"region_name" json:"region_name,computed"`
	RemindExpirationMessage    timetypes.RFC3339                                                              `tfsdk:"remind_expiration_message" json:"remind_expiration_message,computed" format:"date"`
	Resources                  customfield.NestedObjectList[CloudBillingReservationsResourcesDataSourceModel] `tfsdk:"resources" json:"resources,computed"`
	Status                     types.String                                                                   `tfsdk:"status" json:"status,computed"`
	UserStatus                 types.String                                                                   `tfsdk:"user_status" json:"user_status,computed"`
}

type CloudBillingReservationsAmountPricesDataSourceModel struct {
	CommitPricePerMonth     types.String `tfsdk:"commit_price_per_month" json:"commit_price_per_month,computed"`
	CommitPricePerUnit      types.String `tfsdk:"commit_price_per_unit" json:"commit_price_per_unit,computed"`
	CommitPriceTotal        types.String `tfsdk:"commit_price_total" json:"commit_price_total,computed"`
	CurrencyCode            types.String `tfsdk:"currency_code" json:"currency_code,computed"`
	OvercommitPricePerMonth types.String `tfsdk:"overcommit_price_per_month" json:"overcommit_price_per_month,computed"`
	OvercommitPricePerUnit  types.String `tfsdk:"overcommit_price_per_unit" json:"overcommit_price_per_unit,computed"`
	OvercommitPriceTotal    types.String `tfsdk:"overcommit_price_total" json:"overcommit_price_total,computed"`
}

type CloudBillingReservationsResourcesDataSourceModel struct {
	ActivityPeriod              types.String `tfsdk:"activity_period" json:"activity_period,computed"`
	ActivityPeriodLength        types.Int64  `tfsdk:"activity_period_length" json:"activity_period_length,computed"`
	BillingPlanItemID           types.Int64  `tfsdk:"billing_plan_item_id" json:"billing_plan_item_id,computed"`
	CommitPricePerMonth         types.String `tfsdk:"commit_price_per_month" json:"commit_price_per_month,computed"`
	CommitPricePerUnit          types.String `tfsdk:"commit_price_per_unit" json:"commit_price_per_unit,computed"`
	CommitPriceTotal            types.String `tfsdk:"commit_price_total" json:"commit_price_total,computed"`
	OvercommitBillingPlanItemID types.Int64  `tfsdk:"overcommit_billing_plan_item_id" json:"overcommit_billing_plan_item_id,computed"`
	OvercommitPricePerMonth     types.String `tfsdk:"overcommit_price_per_month" json:"overcommit_price_per_month,computed"`
	OvercommitPricePerUnit      types.String `tfsdk:"overcommit_price_per_unit" json:"overcommit_price_per_unit,computed"`
	OvercommitPriceTotal        types.String `tfsdk:"overcommit_price_total" json:"overcommit_price_total,computed"`
	ResourceCount               types.Int64  `tfsdk:"resource_count" json:"resource_count,computed"`
	ResourceName                types.String `tfsdk:"resource_name" json:"resource_name,computed"`
	ResourceType                types.String `tfsdk:"resource_type" json:"resource_type,computed"`
	UnitName                    types.String `tfsdk:"unit_name" json:"unit_name,computed"`
	UnitSizeMonth               types.String `tfsdk:"unit_size_month" json:"unit_size_month,computed"`
	UnitSizeTotal               types.String `tfsdk:"unit_size_total" json:"unit_size_total,computed"`
	CPU                         types.String `tfsdk:"cpu" json:"cpu,computed"`
	Disk                        types.String `tfsdk:"disk" json:"disk,computed"`
	Ram                         types.String `tfsdk:"ram" json:"ram,computed"`
}
