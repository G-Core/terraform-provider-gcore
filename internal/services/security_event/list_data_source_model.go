// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package security_event

import (
	"context"

	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/gcore-go/security"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type SecurityEventsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[SecurityEventsItemsDataSourceModel] `json:"results,computed"`
}

type SecurityEventsDataSourceModel struct {
	AlertType           types.String                                                     `tfsdk:"alert_type" query:"alert_type,optional"`
	TargetedIPAddresses types.String                                                     `tfsdk:"targeted_ip_addresses" query:"targeted_ip_addresses,optional"`
	Limit               types.Int64                                                      `tfsdk:"limit" query:"limit,computed_optional"`
	Ordering            types.String                                                     `tfsdk:"ordering" query:"ordering,computed_optional"`
	MaxItems            types.Int64                                                      `tfsdk:"max_items"`
	Items               customfield.NestedObjectList[SecurityEventsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *SecurityEventsDataSourceModel) toListParams(_ context.Context) (params security.EventListParams, diags diag.Diagnostics) {
	params = security.EventListParams{}

	if !m.AlertType.IsNull() {
		params.AlertType = security.EventListParamsAlertType(m.AlertType.ValueString())
	}
	if !m.Limit.IsNull() {
		params.Limit = param.NewOpt(m.Limit.ValueInt64())
	}
	if !m.Ordering.IsNull() {
		params.Ordering = security.EventListParamsOrdering(m.Ordering.ValueString())
	}
	if !m.TargetedIPAddresses.IsNull() {
		params.TargetedIPAddresses = param.NewOpt(m.TargetedIPAddresses.ValueString())
	}

	return
}

type SecurityEventsItemsDataSourceModel struct {
	ID                         types.String      `tfsdk:"id" json:"id,computed"`
	AlertType                  types.String      `tfsdk:"alert_type" json:"alert_type,computed"`
	AttackPowerBps             types.Float64     `tfsdk:"attack_power_bps" json:"attack_power_bps,computed"`
	AttackPowerPps             types.Float64     `tfsdk:"attack_power_pps" json:"attack_power_pps,computed"`
	AttackStartTime            timetypes.RFC3339 `tfsdk:"attack_start_time" json:"attack_start_time,computed" format:"date-time"`
	ClientID                   types.Int64       `tfsdk:"client_id" json:"client_id,computed"`
	NotificationType           types.String      `tfsdk:"notification_type" json:"notification_type,computed"`
	NumberOfIPInvolvedInAttack types.Int64       `tfsdk:"number_of_ip_involved_in_attack" json:"number_of_ip_involved_in_attack,computed"`
	TargetedIPAddresses        types.String      `tfsdk:"targeted_ip_addresses" json:"targeted_ip_addresses,computed"`
}
