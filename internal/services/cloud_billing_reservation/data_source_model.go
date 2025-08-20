// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_billing_reservation

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudBillingReservationDataSourceModel struct {
	ReservationID              types.Int64                                                                   `tfsdk:"reservation_id" path:"reservation_id,required"`
	ActiveFrom                 timetypes.RFC3339                                                             `tfsdk:"active_from" json:"active_from,computed" format:"date"`
	ActiveTo                   timetypes.RFC3339                                                             `tfsdk:"active_to" json:"active_to,computed" format:"date"`
	ActivityPeriod             types.String                                                                  `tfsdk:"activity_period" json:"activity_period,computed"`
	ActivityPeriodLength       types.Int64                                                                   `tfsdk:"activity_period_length" json:"activity_period_length,computed"`
	BillingPlanID              types.Int64                                                                   `tfsdk:"billing_plan_id" json:"billing_plan_id,computed"`
	CreatedAt                  timetypes.RFC3339                                                             `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	Error                      types.String                                                                  `tfsdk:"error" json:"error,computed"`
	Eta                        timetypes.RFC3339                                                             `tfsdk:"eta" json:"eta,computed" format:"date"`
	ID                         types.Int64                                                                   `tfsdk:"id" json:"id,computed"`
	IsExpirationMessageVisible types.Bool                                                                    `tfsdk:"is_expiration_message_visible" json:"is_expiration_message_visible,computed"`
	Name                       types.String                                                                  `tfsdk:"name" json:"name,computed"`
	RegionID                   types.Int64                                                                   `tfsdk:"region_id" json:"region_id,computed"`
	RegionName                 types.String                                                                  `tfsdk:"region_name" json:"region_name,computed"`
	RemindExpirationMessage    timetypes.RFC3339                                                             `tfsdk:"remind_expiration_message" json:"remind_expiration_message,computed" format:"date"`
	Status                     types.String                                                                  `tfsdk:"status" json:"status,computed"`
	UserStatus                 types.String                                                                  `tfsdk:"user_status" json:"user_status,computed"`
	NextStatuses               customfield.List[types.String]                                                `tfsdk:"next_statuses" json:"next_statuses,computed"`
	AmountPrices               customfield.NestedObject[CloudBillingReservationAmountPricesDataSourceModel]  `tfsdk:"amount_prices" json:"amount_prices,computed"`
	Resources                  customfield.NestedObjectList[CloudBillingReservationResourcesDataSourceModel] `tfsdk:"resources" json:"resources,computed"`
}

type CloudBillingReservationAmountPricesDataSourceModel struct {
	CommitPricePerMonth     types.String `tfsdk:"commit_price_per_month" json:"commit_price_per_month,computed"`
	CommitPricePerUnit      types.String `tfsdk:"commit_price_per_unit" json:"commit_price_per_unit,computed"`
	CommitPriceTotal        types.String `tfsdk:"commit_price_total" json:"commit_price_total,computed"`
	CurrencyCode            types.String `tfsdk:"currency_code" json:"currency_code,computed"`
	OvercommitPricePerMonth types.String `tfsdk:"overcommit_price_per_month" json:"overcommit_price_per_month,computed"`
	OvercommitPricePerUnit  types.String `tfsdk:"overcommit_price_per_unit" json:"overcommit_price_per_unit,computed"`
	OvercommitPriceTotal    types.String `tfsdk:"overcommit_price_total" json:"overcommit_price_total,computed"`
}

type CloudBillingReservationResourcesDataSourceModel struct {
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
