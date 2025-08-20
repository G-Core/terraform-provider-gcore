// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package dns_zone_rrset

import (
	"context"

	"github.com/G-Core/gcore-go/dns"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type DNSZoneRrsetDataSourceModel struct {
	RrsetName       types.String                                                             `tfsdk:"rrset_name" path:"rrsetName,required"`
	RrsetType       types.String                                                             `tfsdk:"rrset_type" path:"rrsetType,required"`
	ZoneName        types.String                                                             `tfsdk:"zone_name" path:"zoneName,required"`
	FilterSetID     types.Int64                                                              `tfsdk:"filter_set_id" json:"filter_set_id,computed"`
	Name            types.String                                                             `tfsdk:"name" json:"name,computed"`
	Ttl             types.Int64                                                              `tfsdk:"ttl" json:"ttl,computed"`
	Type            types.String                                                             `tfsdk:"type" json:"type,computed"`
	UpdatedAt       timetypes.RFC3339                                                        `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	Warning         types.String                                                             `tfsdk:"warning" json:"warning,computed"`
	Meta            customfield.Map[jsontypes.Normalized]                                    `tfsdk:"meta" json:"meta,computed"`
	Pickers         customfield.NestedObjectList[DNSZoneRrsetPickersDataSourceModel]         `tfsdk:"pickers" json:"pickers,computed"`
	ResourceRecords customfield.NestedObjectList[DNSZoneRrsetResourceRecordsDataSourceModel] `tfsdk:"resource_records" json:"resource_records,computed"`
	Warnings        customfield.NestedObjectList[DNSZoneRrsetWarningsDataSourceModel]        `tfsdk:"warnings" json:"warnings,computed"`
}

func (m *DNSZoneRrsetDataSourceModel) toReadParams(_ context.Context) (params dns.ZoneRrsetGetParams, diags diag.Diagnostics) {
	params = dns.ZoneRrsetGetParams{
		ZoneName:  m.ZoneName.ValueString(),
		RrsetName: m.RrsetName.ValueString(),
	}

	return
}

type DNSZoneRrsetPickersDataSourceModel struct {
	Type   types.String `tfsdk:"type" json:"type,computed"`
	Limit  types.Int64  `tfsdk:"limit" json:"limit,computed"`
	Strict types.Bool   `tfsdk:"strict" json:"strict,computed"`
}

type DNSZoneRrsetResourceRecordsDataSourceModel struct {
	Content customfield.List[jsontypes.Normalized] `tfsdk:"content" json:"content,computed"`
	ID      types.Int64                            `tfsdk:"id" json:"id,computed"`
	Enabled types.Bool                             `tfsdk:"enabled" json:"enabled,computed"`
	Meta    customfield.Map[jsontypes.Normalized]  `tfsdk:"meta" json:"meta,computed"`
}

type DNSZoneRrsetWarningsDataSourceModel struct {
	Key     types.String `tfsdk:"key" json:"key,computed"`
	Message types.String `tfsdk:"message" json:"message,computed"`
}
