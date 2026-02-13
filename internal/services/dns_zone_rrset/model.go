// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package dns_zone_rrset

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type DNSZoneRrsetModel struct {
	RrsetName       types.String                                            `tfsdk:"rrset_name" path:"rrsetName,required"`
	RrsetType       types.String                                            `tfsdk:"rrset_type" path:"rrsetType,required"`
	ZoneName        types.String                                            `tfsdk:"zone_name" path:"zoneName,required"`
	ResourceRecords *[]*DNSZoneRrsetResourceRecordsModel                    `tfsdk:"resource_records" json:"resource_records,required"`
	Ttl             types.Int64                                             `tfsdk:"ttl" json:"ttl,optional"`
	Meta            *map[string]customfield.MetaStringValue                 `tfsdk:"meta" json:"meta,optional"`
	Pickers         *[]*DNSZoneRrsetPickersModel                            `tfsdk:"pickers" json:"pickers,optional"`
	FilterSetID     types.Int64                                             `tfsdk:"filter_set_id" json:"filter_set_id,computed"`
	Name            types.String                                            `tfsdk:"name" json:"name,computed"`
	Type            types.String                                            `tfsdk:"type" json:"type,computed"`
	Warning         types.String                                            `tfsdk:"warning" json:"warning,computed"`
	Warnings        customfield.NestedObjectList[DNSZoneRrsetWarningsModel] `tfsdk:"warnings" json:"warnings,computed"`
}

func (m DNSZoneRrsetModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m DNSZoneRrsetModel) MarshalJSONForUpdate(state DNSZoneRrsetModel) (data []byte, err error) {
	return apijson.MarshalForUpdate(m, state)
}

type DNSZoneRrsetResourceRecordsModel struct {
	ID      types.Int64                             `tfsdk:"id" json:"id,computed"`
	Content *[]jsontypes.Normalized                 `tfsdk:"content" json:"content,required"`
	Enabled types.Bool                              `tfsdk:"enabled" json:"enabled,computed_optional"`
	Meta    *map[string]customfield.MetaStringValue `tfsdk:"meta" json:"meta,optional"`
}

type DNSZoneRrsetPickersModel struct {
	Type   types.String `tfsdk:"type" json:"type,required"`
	Limit  types.Int64  `tfsdk:"limit" json:"limit,computed_optional"`
	Strict types.Bool   `tfsdk:"strict" json:"strict,computed_optional"`
}

type DNSZoneRrsetWarningsModel struct {
	Key     types.String `tfsdk:"key" json:"key,computed"`
	Message types.String `tfsdk:"message" json:"message,computed"`
}
