// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package dns_zone

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type DNSZoneDataSourceModel struct {
	Name          types.String                                                 `tfsdk:"name" path:"name,required"`
	Contact       types.String                                                 `tfsdk:"contact" json:"contact,computed"`
	DnssecEnabled types.Bool                                                   `tfsdk:"dnssec_enabled" json:"dnssec_enabled,computed"`
	Enabled       types.Bool                                                   `tfsdk:"enabled" json:"enabled,computed"`
	Expiry        types.Int64                                                  `tfsdk:"expiry" json:"expiry,computed"`
	ID            types.Int64                                                  `tfsdk:"id" json:"id,computed"`
	NxTtl         types.Int64                                                  `tfsdk:"nx_ttl" json:"nx_ttl,computed"`
	PrimaryServer types.String                                                 `tfsdk:"primary_server" json:"primary_server,computed"`
	Refresh       types.Int64                                                  `tfsdk:"refresh" json:"refresh,computed"`
	Retry         types.Int64                                                  `tfsdk:"retry" json:"retry,computed"`
	Serial        types.Int64                                                  `tfsdk:"serial" json:"serial,computed"`
	Status        types.String                                                 `tfsdk:"status" json:"status,computed"`
	Meta          customfield.Map[jsontypes.Normalized]                        `tfsdk:"meta" json:"meta,computed"`
	Records       customfield.NestedObjectList[DNSZoneRecordsDataSourceModel]  `tfsdk:"records" json:"records,computed"`
	RrsetsAmount  customfield.NestedObject[DNSZoneRrsetsAmountDataSourceModel] `tfsdk:"rrsets_amount" json:"rrsets_amount,computed"`
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
