// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package dns_zone

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type DNSZoneModel struct {
	ID            types.String                                       `tfsdk:"id" json:"-,computed"`
	Name          types.String                                       `tfsdk:"name" json:"name,required"`
	Contact       types.String                                       `tfsdk:"contact" json:"contact,computed_optional"`
	Expiry        types.Int64                                        `tfsdk:"expiry" json:"expiry,computed_optional"`
	NxTtl         types.Int64                                        `tfsdk:"nx_ttl" json:"nx_ttl,computed_optional"`
	PrimaryServer types.String                                       `tfsdk:"primary_server" json:"primary_server,computed_optional"`
	Refresh       types.Int64                                        `tfsdk:"refresh" json:"refresh,computed_optional"`
	Retry         types.Int64                                        `tfsdk:"retry" json:"retry,computed_optional"`
	Serial        types.Int64                                        `tfsdk:"serial" json:"serial,computed_optional"`
	Meta          *map[string]customfield.MetaStringValue             `tfsdk:"meta" json:"meta,optional"`
	Enabled       types.Bool                                         `tfsdk:"enabled" json:"enabled,computed_optional"`
	DnssecEnabled types.Bool                                         `tfsdk:"dnssec_enabled" json:"dnssec_enabled,computed"`
	Status        types.String                                       `tfsdk:"status" json:"status,computed"`
	Warnings      customfield.List[types.String]                     `tfsdk:"warnings" json:"warnings,computed"`
	Records       customfield.NestedObjectList[DNSZoneRecordsModel]  `tfsdk:"records" json:"records,computed"`
	RrsetsAmount  customfield.NestedObject[DNSZoneRrsetsAmountModel] `tfsdk:"rrsets_amount" json:"rrsets_amount,computed"`
}

func (m DNSZoneModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m DNSZoneModel) MarshalJSONForUpdate(state DNSZoneModel) (data []byte, err error) {
	return apijson.MarshalForUpdate(m, state)
}

type DNSZoneRecordsModel struct {
	Name         types.String                   `tfsdk:"name" json:"name,computed"`
	ShortAnswers customfield.List[types.String] `tfsdk:"short_answers" json:"short_answers,computed"`
	Ttl          types.Int64                    `tfsdk:"ttl" json:"ttl,computed"`
	Type         types.String                   `tfsdk:"type" json:"type,computed"`
}

type DNSZoneRrsetsAmountModel struct {
	Dynamic customfield.NestedObject[DNSZoneRrsetsAmountDynamicModel] `tfsdk:"dynamic" json:"dynamic,computed"`
	Static  types.Int64                                               `tfsdk:"static" json:"static,computed"`
	Total   types.Int64                                               `tfsdk:"total" json:"total,computed"`
}

type DNSZoneRrsetsAmountDynamicModel struct {
	Healthcheck types.Int64 `tfsdk:"healthcheck" json:"healthcheck,computed"`
	Total       types.Int64 `tfsdk:"total" json:"total,computed"`
}
