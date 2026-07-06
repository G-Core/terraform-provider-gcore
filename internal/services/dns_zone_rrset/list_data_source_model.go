// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package dns_zone_rrset

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

type DNSZoneRrsetsRrsetsListDataSourceEnvelope struct {
	Rrsets customfield.NestedObjectList[DNSZoneRrsetsItemsDataSourceModel] `json:"rrsets,computed"`
}

type DNSZoneRrsetsDataSourceModel struct {
	ZoneName       types.String                                                    `tfsdk:"zone_name" path:"zoneName,required"`
	OrderBy        types.String                                                    `tfsdk:"order_by" query:"order_by,optional"`
	OrderDirection types.String                                                    `tfsdk:"order_direction" query:"order_direction,optional"`
	MaxItems       types.Int64                                                     `tfsdk:"max_items"`
	Items          customfield.NestedObjectList[DNSZoneRrsetsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *DNSZoneRrsetsDataSourceModel) toListParams(_ context.Context) (params dns.ZoneRrsetListParams, diags diag.Diagnostics) {
	params = dns.ZoneRrsetListParams{}

	if !m.OrderBy.IsNull() {
		params.OrderBy = param.NewOpt(m.OrderBy.ValueString())
	}
	if !m.OrderDirection.IsNull() {
		params.OrderDirection = dns.ZoneRrsetListParamsOrderDirection(m.OrderDirection.ValueString())
	}

	return
}

type DNSZoneRrsetsItemsDataSourceModel struct {
	Name            types.String                                                              `tfsdk:"name" json:"name,computed"`
	ResourceRecords customfield.NestedObjectList[DNSZoneRrsetsResourceRecordsDataSourceModel] `tfsdk:"resource_records" json:"resource_records,computed"`
	Type            types.String                                                              `tfsdk:"type" json:"type,computed"`
	FilterSetID     types.Int64                                                               `tfsdk:"filter_set_id" json:"filter_set_id,computed"`
	Meta            customfield.Map[jsontypes.Normalized]                                     `tfsdk:"meta" json:"meta,computed"`
	Pickers         customfield.NestedObjectList[DNSZoneRrsetsPickersDataSourceModel]         `tfsdk:"pickers" json:"pickers,computed"`
	Ttl             types.Int64                                                               `tfsdk:"ttl" json:"ttl,computed"`
	UpdatedAt       timetypes.RFC3339                                                         `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	Warning         types.String                                                              `tfsdk:"warning" json:"warning,computed"`
	Warnings        customfield.NestedObjectList[DNSZoneRrsetsWarningsDataSourceModel]        `tfsdk:"warnings" json:"warnings,computed"`
}

type DNSZoneRrsetsResourceRecordsDataSourceModel struct {
	Content customfield.List[jsontypes.Normalized] `tfsdk:"content" json:"content,computed"`
	ID      types.Int64                            `tfsdk:"id" json:"id,computed"`
	Enabled types.Bool                             `tfsdk:"enabled" json:"enabled,computed"`
	Meta    customfield.Map[jsontypes.Normalized]  `tfsdk:"meta" json:"meta,computed"`
}

type DNSZoneRrsetsPickersDataSourceModel struct {
	Type   types.String `tfsdk:"type" json:"type,computed"`
	Limit  types.Int64  `tfsdk:"limit" json:"limit,computed"`
	Strict types.Bool   `tfsdk:"strict" json:"strict,computed"`
}

type DNSZoneRrsetsWarningsDataSourceModel struct {
	Key     types.String `tfsdk:"key" json:"key,computed"`
	Message types.String `tfsdk:"message" json:"message,computed"`
}
