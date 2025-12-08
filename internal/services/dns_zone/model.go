// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package dns_zone

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type DNSZoneModel struct {
	ID            types.Int64                                `tfsdk:"id" json:"id,computed"`
	Name          types.String                               `tfsdk:"name" json:"name,required,no_refresh"`
	Contact       types.String                               `tfsdk:"contact" json:"contact,optional,no_refresh"`
	Expiry        types.Int64                                `tfsdk:"expiry" json:"expiry,optional,no_refresh"`
	NxTtl         types.Int64                                `tfsdk:"nx_ttl" json:"nx_ttl,optional,no_refresh"`
	PrimaryServer types.String                               `tfsdk:"primary_server" json:"primary_server,optional,no_refresh"`
	Refresh       types.Int64                                `tfsdk:"refresh" json:"refresh,optional,no_refresh"`
	Retry         types.Int64                                `tfsdk:"retry" json:"retry,optional,no_refresh"`
	Serial        types.Int64                                `tfsdk:"serial" json:"serial,optional,no_refresh"`
	Meta          *map[string]jsontypes.Normalized           `tfsdk:"meta" json:"meta,optional,no_refresh"`
	Enabled       types.Bool                                 `tfsdk:"enabled" json:"enabled,computed_optional,no_refresh"`
	Warnings      customfield.List[types.String]             `tfsdk:"warnings" json:"warnings,computed,no_refresh"`
	Zone          customfield.NestedObject[DNSZoneZoneModel] `tfsdk:"zone" json:"Zone,computed"`
}

func (m DNSZoneModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m DNSZoneModel) MarshalJSONForUpdate(state DNSZoneModel) (data []byte, err error) {
	return apijson.MarshalForUpdate(m, state)
}

type DNSZoneZoneModel struct {
	ID            types.Int64                                            `tfsdk:"id" json:"id,computed"`
	ClientID      types.Int64                                            `tfsdk:"client_id" json:"client_id,computed"`
	Contact       types.String                                           `tfsdk:"contact" json:"contact,computed"`
	DnssecEnabled types.Bool                                             `tfsdk:"dnssec_enabled" json:"dnssec_enabled,computed"`
	Expiry        types.Int64                                            `tfsdk:"expiry" json:"expiry,computed"`
	Meta          jsontypes.Normalized                                   `tfsdk:"meta" json:"meta,computed"`
	Name          types.String                                           `tfsdk:"name" json:"name,computed"`
	NxTtl         types.Int64                                            `tfsdk:"nx_ttl" json:"nx_ttl,computed"`
	PrimaryServer types.String                                           `tfsdk:"primary_server" json:"primary_server,computed"`
	Records       customfield.NestedObjectList[DNSZoneZoneRecordsModel]  `tfsdk:"records" json:"records,computed"`
	Refresh       types.Int64                                            `tfsdk:"refresh" json:"refresh,computed"`
	Retry         types.Int64                                            `tfsdk:"retry" json:"retry,computed"`
	RrsetsAmount  customfield.NestedObject[DNSZoneZoneRrsetsAmountModel] `tfsdk:"rrsets_amount" json:"rrsets_amount,computed"`
	Serial        types.Int64                                            `tfsdk:"serial" json:"serial,computed"`
	Status        types.String                                           `tfsdk:"status" json:"status,computed"`
}

type DNSZoneZoneRecordsModel struct {
	Name         types.String                   `tfsdk:"name" json:"name,computed"`
	ShortAnswers customfield.List[types.String] `tfsdk:"short_answers" json:"short_answers,computed"`
	Ttl          types.Int64                    `tfsdk:"ttl" json:"ttl,computed"`
	Type         types.String                   `tfsdk:"type" json:"type,computed"`
}

type DNSZoneZoneRrsetsAmountModel struct {
	Dynamic customfield.NestedObject[DNSZoneZoneRrsetsAmountDynamicModel] `tfsdk:"dynamic" json:"dynamic,computed"`
	Static  types.Int64                                                   `tfsdk:"static" json:"static,computed"`
	Total   types.Int64                                                   `tfsdk:"total" json:"total,computed"`
}

type DNSZoneZoneRrsetsAmountDynamicModel struct {
	Healthcheck types.Int64 `tfsdk:"healthcheck" json:"healthcheck,computed"`
	Total       types.Int64 `tfsdk:"total" json:"total,computed"`
}
