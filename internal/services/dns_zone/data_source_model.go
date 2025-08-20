// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package dns_zone

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type DNSZoneDataSourceModel struct {
	Name types.String                                         `tfsdk:"name" path:"name,required"`
	Zone customfield.NestedObject[DNSZoneZoneDataSourceModel] `tfsdk:"zone" json:"Zone,computed"`
}

type DNSZoneZoneDataSourceModel struct {
	ID            types.Int64                                                      `tfsdk:"id" json:"id,computed"`
	Contact       types.String                                                     `tfsdk:"contact" json:"contact,computed"`
	DnssecEnabled types.Bool                                                       `tfsdk:"dnssec_enabled" json:"dnssec_enabled,computed"`
	Expiry        types.Int64                                                      `tfsdk:"expiry" json:"expiry,computed"`
	Meta          jsontypes.Normalized                                             `tfsdk:"meta" json:"meta,computed"`
	Name          types.String                                                     `tfsdk:"name" json:"name,computed"`
	NxTtl         types.Int64                                                      `tfsdk:"nx_ttl" json:"nx_ttl,computed"`
	PrimaryServer types.String                                                     `tfsdk:"primary_server" json:"primary_server,computed"`
	Records       customfield.NestedObjectList[DNSZoneZoneRecordsDataSourceModel]  `tfsdk:"records" json:"records,computed"`
	Refresh       types.Int64                                                      `tfsdk:"refresh" json:"refresh,computed"`
	Retry         types.Int64                                                      `tfsdk:"retry" json:"retry,computed"`
	RrsetsAmount  customfield.NestedObject[DNSZoneZoneRrsetsAmountDataSourceModel] `tfsdk:"rrsets_amount" json:"rrsets_amount,computed"`
	Serial        types.Int64                                                      `tfsdk:"serial" json:"serial,computed"`
	Status        types.String                                                     `tfsdk:"status" json:"status,computed"`
}

type DNSZoneZoneRecordsDataSourceModel struct {
	Name         types.String                   `tfsdk:"name" json:"name,computed"`
	ShortAnswers customfield.List[types.String] `tfsdk:"short_answers" json:"short_answers,computed"`
	Ttl          types.Int64                    `tfsdk:"ttl" json:"ttl,computed"`
	Type         types.String                   `tfsdk:"type" json:"type,computed"`
}

type DNSZoneZoneRrsetsAmountDataSourceModel struct {
	Dynamic customfield.NestedObject[DNSZoneZoneRrsetsAmountDynamicDataSourceModel] `tfsdk:"dynamic" json:"dynamic,computed"`
	Static  types.Int64                                                             `tfsdk:"static" json:"static,computed"`
	Total   types.Int64                                                             `tfsdk:"total" json:"total,computed"`
}

type DNSZoneZoneRrsetsAmountDynamicDataSourceModel struct {
	Healthcheck types.Int64 `tfsdk:"healthcheck" json:"healthcheck,computed"`
	Total       types.Int64 `tfsdk:"total" json:"total,computed"`
}
