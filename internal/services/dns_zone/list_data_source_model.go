// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package dns_zone

import (
	"context"

	"github.com/G-Core/gcore-go/dns"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DNSZonesZonesListDataSourceEnvelope struct {
	Zones customfield.NestedObjectList[DNSZonesItemsDataSourceModel] `json:"zones,computed"`
}

type DNSZonesDataSourceModel struct {
	CaseSensitive  types.Bool                                                 `tfsdk:"case_sensitive" query:"case_sensitive,optional"`
	Dynamic        types.Bool                                                 `tfsdk:"dynamic" query:"dynamic,optional"`
	Enabled        types.Bool                                                 `tfsdk:"enabled" query:"enabled,optional"`
	ExactMatch     types.Bool                                                 `tfsdk:"exact_match" query:"exact_match,optional"`
	Healthcheck    types.Bool                                                 `tfsdk:"healthcheck" query:"healthcheck,optional"`
	OrderBy        types.String                                               `tfsdk:"order_by" query:"order_by,optional"`
	OrderDirection types.String                                               `tfsdk:"order_direction" query:"order_direction,optional"`
	Status         types.String                                               `tfsdk:"status" query:"status,optional"`
	UpdatedAtFrom  timetypes.RFC3339                                          `tfsdk:"updated_at_from" query:"updated_at_from,optional" format:"date-time"`
	UpdatedAtTo    timetypes.RFC3339                                          `tfsdk:"updated_at_to" query:"updated_at_to,optional" format:"date-time"`
	ClientID       *[]types.Int64                                             `tfsdk:"client_id" query:"client_id,optional"`
	IamResellerID  *[]types.Int64                                             `tfsdk:"iam_reseller_id" query:"iam_reseller_id,optional"`
	ID             *[]types.Int64                                             `tfsdk:"id" query:"id,optional"`
	Name           *[]types.String                                            `tfsdk:"name" query:"name,optional"`
	ResellerID     *[]types.Int64                                             `tfsdk:"reseller_id" query:"reseller_id,optional"`
	MaxItems       types.Int64                                                `tfsdk:"max_items"`
	Items          customfield.NestedObjectList[DNSZonesItemsDataSourceModel] `tfsdk:"items"`
}

func (m *DNSZonesDataSourceModel) toListParams(_ context.Context) (params dns.ZoneListParams, diags diag.Diagnostics) {
	mID := []int64{}
	if m.ID != nil {
		for _, item := range *m.ID {
			mID = append(mID, item.ValueInt64())
		}
	}
	mClientID := []int64{}
	if m.ClientID != nil {
		for _, item := range *m.ClientID {
			mClientID = append(mClientID, item.ValueInt64())
		}
	}
	mIamResellerID := []int64{}
	if m.IamResellerID != nil {
		for _, item := range *m.IamResellerID {
			mIamResellerID = append(mIamResellerID, item.ValueInt64())
		}
	}
	mName := []string{}
	if m.Name != nil {
		for _, item := range *m.Name {
			mName = append(mName, item.ValueString())
		}
	}
	mResellerID := []int64{}
	if m.ResellerID != nil {
		for _, item := range *m.ResellerID {
			mResellerID = append(mResellerID, item.ValueInt64())
		}
	}

	params = dns.ZoneListParams{
		ID:            mID,
		ClientID:      mClientID,
		IamResellerID: mIamResellerID,
		Name:          mName,
		ResellerID:    mResellerID,
	}

	if !m.CaseSensitive.IsNull() {
		params.CaseSensitive = param.NewOpt(m.CaseSensitive.ValueBool())
	}
	if !m.Dynamic.IsNull() {
		params.Dynamic = param.NewOpt(m.Dynamic.ValueBool())
	}
	if !m.Enabled.IsNull() {
		params.Enabled = param.NewOpt(m.Enabled.ValueBool())
	}
	if !m.ExactMatch.IsNull() {
		params.ExactMatch = param.NewOpt(m.ExactMatch.ValueBool())
	}
	if !m.Healthcheck.IsNull() {
		params.Healthcheck = param.NewOpt(m.Healthcheck.ValueBool())
	}
	if !m.OrderBy.IsNull() {
		params.OrderBy = param.NewOpt(m.OrderBy.ValueString())
	}
	if !m.OrderDirection.IsNull() {
		params.OrderDirection = dns.ZoneListParamsOrderDirection(m.OrderDirection.ValueString())
	}
	if !m.Status.IsNull() {
		params.Status = param.NewOpt(m.Status.ValueString())
	}
	if !m.UpdatedAtFrom.IsNull() {
		mUpdatedAtFrom, errs := m.UpdatedAtFrom.ValueRFC3339Time()
		diags.Append(errs...)
		params.UpdatedAtFrom = param.NewOpt(mUpdatedAtFrom)
	}
	if !m.UpdatedAtTo.IsNull() {
		mUpdatedAtTo, errs := m.UpdatedAtTo.ValueRFC3339Time()
		diags.Append(errs...)
		params.UpdatedAtTo = param.NewOpt(mUpdatedAtTo)
	}

	return
}

type DNSZonesItemsDataSourceModel struct {
	ID                     types.Int64                                                   `tfsdk:"id" json:"id,computed"`
	ClientID               types.Int64                                                   `tfsdk:"client_id" json:"client_id,computed"`
	Contact                types.String                                                  `tfsdk:"contact" json:"contact,computed"`
	DnssecEnabled          types.Bool                                                    `tfsdk:"dnssec_enabled" json:"dnssec_enabled,computed"`
	DnssecStatus           types.String                                                  `tfsdk:"dnssec_status" json:"dnssec_status,computed"`
	DnssecStatusModifiedOn types.String                                                  `tfsdk:"dnssec_status_modified_on" json:"dnssec_status_modified_on,computed"`
	Enabled                types.Bool                                                    `tfsdk:"enabled" json:"enabled,computed"`
	Expiry                 types.Int64                                                   `tfsdk:"expiry" json:"expiry,computed"`
	Meta                   customfield.Map[jsontypes.Normalized]                         `tfsdk:"meta" json:"meta,computed"`
	Name                   types.String                                                  `tfsdk:"name" json:"name,computed"`
	NxTtl                  types.Int64                                                   `tfsdk:"nx_ttl" json:"nx_ttl,computed"`
	PrimaryServer          types.String                                                  `tfsdk:"primary_server" json:"primary_server,computed"`
	Records                customfield.NestedObjectList[DNSZonesRecordsDataSourceModel]  `tfsdk:"records" json:"records,computed"`
	Refresh                types.Int64                                                   `tfsdk:"refresh" json:"refresh,computed"`
	Retry                  types.Int64                                                   `tfsdk:"retry" json:"retry,computed"`
	RrsetsAmount           customfield.NestedObject[DNSZonesRrsetsAmountDataSourceModel] `tfsdk:"rrsets_amount" json:"rrsets_amount,computed"`
	Serial                 types.Int64                                                   `tfsdk:"serial" json:"serial,computed"`
	Status                 types.String                                                  `tfsdk:"status" json:"status,computed"`
}

type DNSZonesRecordsDataSourceModel struct {
	Name         types.String                   `tfsdk:"name" json:"name,computed"`
	ShortAnswers customfield.List[types.String] `tfsdk:"short_answers" json:"short_answers,computed"`
	Ttl          types.Int64                    `tfsdk:"ttl" json:"ttl,computed"`
	Type         types.String                   `tfsdk:"type" json:"type,computed"`
}

type DNSZonesRrsetsAmountDataSourceModel struct {
	Dynamic customfield.NestedObject[DNSZonesRrsetsAmountDynamicDataSourceModel] `tfsdk:"dynamic" json:"dynamic,computed"`
	Static  types.Int64                                                          `tfsdk:"static" json:"static,computed"`
	Total   types.Int64                                                          `tfsdk:"total" json:"total,computed"`
}

type DNSZonesRrsetsAmountDynamicDataSourceModel struct {
	Healthcheck types.Int64 `tfsdk:"healthcheck" json:"healthcheck,computed"`
	Total       types.Int64 `tfsdk:"total" json:"total,computed"`
}
