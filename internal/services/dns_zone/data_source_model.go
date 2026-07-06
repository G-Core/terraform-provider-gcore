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

type DNSZoneDataSourceModel struct {
	ID                     types.String                                                 `tfsdk:"id" path:"name,computed"`
	Name                   types.String                                                 `tfsdk:"name" path:"name,computed_optional"`
	Contact                types.String                                                 `tfsdk:"contact" json:"contact,computed"`
	DnssecEnabled          types.Bool                                                   `tfsdk:"dnssec_enabled" json:"dnssec_enabled,computed"`
	DnssecStatus           types.String                                                 `tfsdk:"dnssec_status" json:"dnssec_status,computed"`
	DnssecStatusModifiedOn types.String                                                 `tfsdk:"dnssec_status_modified_on" json:"dnssec_status_modified_on,computed"`
	Enabled                types.Bool                                                   `tfsdk:"enabled" json:"enabled,computed"`
	Expiry                 types.Int64                                                  `tfsdk:"expiry" json:"expiry,computed"`
	NxTtl                  types.Int64                                                  `tfsdk:"nx_ttl" json:"nx_ttl,computed"`
	PrimaryServer          types.String                                                 `tfsdk:"primary_server" json:"primary_server,computed"`
	Refresh                types.Int64                                                  `tfsdk:"refresh" json:"refresh,computed"`
	Retry                  types.Int64                                                  `tfsdk:"retry" json:"retry,computed"`
	Serial                 types.Int64                                                  `tfsdk:"serial" json:"serial,computed"`
	Status                 types.String                                                 `tfsdk:"status" json:"status,computed"`
	Meta                   customfield.Map[jsontypes.Normalized]                        `tfsdk:"meta" json:"meta,computed"`
	Records                customfield.NestedObjectList[DNSZoneRecordsDataSourceModel]  `tfsdk:"records" json:"records,computed"`
	RrsetsAmount           customfield.NestedObject[DNSZoneRrsetsAmountDataSourceModel] `tfsdk:"rrsets_amount" json:"rrsets_amount,computed"`
	FindOneBy              *DNSZoneFindOneByDataSourceModel                             `tfsdk:"find_one_by"`
}

func (m *DNSZoneDataSourceModel) toListParams(_ context.Context) (params dns.ZoneListParams, diags diag.Diagnostics) {
	mFindOneByID := []int64{}
	if m.FindOneBy.ID != nil {
		for _, item := range *m.FindOneBy.ID {
			mFindOneByID = append(mFindOneByID, item.ValueInt64())
		}
	}
	mFindOneByClientID := []int64{}
	if m.FindOneBy.ClientID != nil {
		for _, item := range *m.FindOneBy.ClientID {
			mFindOneByClientID = append(mFindOneByClientID, item.ValueInt64())
		}
	}
	mFindOneByIamResellerID := []int64{}
	if m.FindOneBy.IamResellerID != nil {
		for _, item := range *m.FindOneBy.IamResellerID {
			mFindOneByIamResellerID = append(mFindOneByIamResellerID, item.ValueInt64())
		}
	}
	mFindOneByName := []string{}
	if m.FindOneBy.Name != nil {
		for _, item := range *m.FindOneBy.Name {
			mFindOneByName = append(mFindOneByName, item.ValueString())
		}
	}
	mFindOneByResellerID := []int64{}
	if m.FindOneBy.ResellerID != nil {
		for _, item := range *m.FindOneBy.ResellerID {
			mFindOneByResellerID = append(mFindOneByResellerID, item.ValueInt64())
		}
	}
	mFindOneByUpdatedAtFrom, errs := m.FindOneBy.UpdatedAtFrom.ValueRFC3339Time()
	diags.Append(errs...)
	mFindOneByUpdatedAtTo, errs := m.FindOneBy.UpdatedAtTo.ValueRFC3339Time()
	diags.Append(errs...)

	params = dns.ZoneListParams{
		ID:            mFindOneByID,
		ClientID:      mFindOneByClientID,
		IamResellerID: mFindOneByIamResellerID,
		Name:          mFindOneByName,
		ResellerID:    mFindOneByResellerID,
	}

	if !m.FindOneBy.CaseSensitive.IsNull() {
		params.CaseSensitive = param.NewOpt(m.FindOneBy.CaseSensitive.ValueBool())
	}
	if !m.FindOneBy.Dynamic.IsNull() {
		params.Dynamic = param.NewOpt(m.FindOneBy.Dynamic.ValueBool())
	}
	if !m.FindOneBy.Enabled.IsNull() {
		params.Enabled = param.NewOpt(m.FindOneBy.Enabled.ValueBool())
	}
	if !m.FindOneBy.ExactMatch.IsNull() {
		params.ExactMatch = param.NewOpt(m.FindOneBy.ExactMatch.ValueBool())
	}
	if !m.FindOneBy.Healthcheck.IsNull() {
		params.Healthcheck = param.NewOpt(m.FindOneBy.Healthcheck.ValueBool())
	}
	if !m.FindOneBy.OrderBy.IsNull() {
		params.OrderBy = param.NewOpt(m.FindOneBy.OrderBy.ValueString())
	}
	if !m.FindOneBy.OrderDirection.IsNull() {
		params.OrderDirection = dns.ZoneListParamsOrderDirection(m.FindOneBy.OrderDirection.ValueString())
	}
	if !m.FindOneBy.Status.IsNull() {
		params.Status = param.NewOpt(m.FindOneBy.Status.ValueString())
	}
	if !m.FindOneBy.UpdatedAtFrom.IsNull() {
		params.UpdatedAtFrom = param.NewOpt(mFindOneByUpdatedAtFrom)
	}
	if !m.FindOneBy.UpdatedAtTo.IsNull() {
		params.UpdatedAtTo = param.NewOpt(mFindOneByUpdatedAtTo)
	}

	return
}

type DNSZoneRecordsDataSourceModel struct {
	Name         types.String                   `tfsdk:"name" json:"name,computed"`
	ShortAnswers customfield.List[types.String] `tfsdk:"short_answers" json:"short_answers,computed"`
	Ttl          types.Int64                    `tfsdk:"ttl" json:"ttl,computed"`
	Type         types.String                   `tfsdk:"type" json:"type,computed"`
}

type DNSZoneRrsetsAmountDataSourceModel struct {
	Dynamic customfield.NestedObject[DNSZoneRrsetsAmountDynamicDataSourceModel] `tfsdk:"dynamic" json:"dynamic,computed"`
	Static  types.Int64                                                         `tfsdk:"static" json:"static,computed"`
	Total   types.Int64                                                         `tfsdk:"total" json:"total,computed"`
}

type DNSZoneRrsetsAmountDynamicDataSourceModel struct {
	Healthcheck types.Int64 `tfsdk:"healthcheck" json:"healthcheck,computed"`
	Total       types.Int64 `tfsdk:"total" json:"total,computed"`
}

type DNSZoneFindOneByDataSourceModel struct {
	ID             *[]types.Int64    `tfsdk:"id" query:"id,optional"`
	CaseSensitive  types.Bool        `tfsdk:"case_sensitive" query:"case_sensitive,optional"`
	ClientID       *[]types.Int64    `tfsdk:"client_id" query:"client_id,optional"`
	Dynamic        types.Bool        `tfsdk:"dynamic" query:"dynamic,optional"`
	Enabled        types.Bool        `tfsdk:"enabled" query:"enabled,optional"`
	ExactMatch     types.Bool        `tfsdk:"exact_match" query:"exact_match,optional"`
	Healthcheck    types.Bool        `tfsdk:"healthcheck" query:"healthcheck,optional"`
	IamResellerID  *[]types.Int64    `tfsdk:"iam_reseller_id" query:"iam_reseller_id,optional"`
	Name           *[]types.String   `tfsdk:"name" query:"name,optional"`
	OrderBy        types.String      `tfsdk:"order_by" query:"order_by,optional"`
	OrderDirection types.String      `tfsdk:"order_direction" query:"order_direction,optional"`
	ResellerID     *[]types.Int64    `tfsdk:"reseller_id" query:"reseller_id,optional"`
	Status         types.String      `tfsdk:"status" query:"status,optional"`
	UpdatedAtFrom  timetypes.RFC3339 `tfsdk:"updated_at_from" query:"updated_at_from,optional" format:"date-time"`
	UpdatedAtTo    timetypes.RFC3339 `tfsdk:"updated_at_to" query:"updated_at_to,optional" format:"date-time"`
}
