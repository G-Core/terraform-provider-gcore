// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package dns_zone

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigValidators = (*DNSZoneDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "DNS zones are authoritative containers for domain name records, with support for DNSSEC and SOA configuration.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"contact": schema.StringAttribute{
				Description: "email address of the administrator responsible for this zone",
				Computed:    true,
			},
			"dnssec_enabled": schema.BoolAttribute{
				Description: "describe dnssec status\ntrue means dnssec is enabled for the zone\nfalse means dnssec is disabled for the zone",
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Computed: true,
			},
			"expiry": schema.Int64Attribute{
				Description: "number of seconds after which secondary name servers should stop answering request for this zone",
				Computed:    true,
			},
			"id": schema.Int64Attribute{
				Description: "ID of zone.\nThis field usually is omitted in response and available only in\ncase of getting deleted zones by admin.",
				Computed:    true,
			},
			"nx_ttl": schema.Int64Attribute{
				Description: "Time To Live of cache",
				Computed:    true,
			},
			"primary_server": schema.StringAttribute{
				Description: "primary master name server for zone",
				Computed:    true,
			},
			"refresh": schema.Int64Attribute{
				Description: "number of seconds after which secondary name servers should query the master for the SOA record, to detect zone changes.",
				Computed:    true,
			},
			"retry": schema.Int64Attribute{
				Description: "number of seconds after which secondary name servers should retry to request the serial number",
				Computed:    true,
			},
			"serial": schema.Int64Attribute{
				Description: "Serial number for this zone or Timestamp of zone modification moment.\nIf a secondary name server slaved to this one observes an increase in this number,\nthe slave will assume that the zone has been updated and initiate a zone transfer.",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
			"meta": schema.MapAttribute{
				Description: "arbitrarily data of zone in json format",
				Computed:    true,
				CustomType:  customfield.NewMapType[jsontypes.Normalized](ctx),
				ElementType: jsontypes.NormalizedType{},
			},
			"records": schema.ListNestedAttribute{
				Computed:   true,
				CustomType: customfield.NewNestedObjectListType[DNSZoneRecordsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"short_answers": schema.ListAttribute{
							Computed:    true,
							CustomType:  customfield.NewListType[types.String](ctx),
							ElementType: types.StringType,
						},
						"ttl": schema.Int64Attribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"rrsets_amount": schema.SingleNestedAttribute{
				Computed:   true,
				CustomType: customfield.NewNestedObjectType[DNSZoneRrsetsAmountDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"dynamic": schema.SingleNestedAttribute{
						Description: "Amount of dynamic RRsets in zone",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[DNSZoneRrsetsAmountDynamicDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"healthcheck": schema.Int64Attribute{
								Description: "Amount of RRsets with enabled healthchecks",
								Computed:    true,
							},
							"total": schema.Int64Attribute{
								Description: "Total amount of dynamic RRsets in zone",
								Computed:    true,
							},
						},
					},
					"static": schema.Int64Attribute{
						Description: "Amount of static RRsets in zone",
						Computed:    true,
					},
					"total": schema.Int64Attribute{
						Description: "Total amount of RRsets in zone",
						Computed:    true,
					},
				},
			},
		},
	}
}

func (d *DNSZoneDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *DNSZoneDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
