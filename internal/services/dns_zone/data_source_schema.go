// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package dns_zone

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigValidators = (*DNSZoneDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "DNS zones are authoritative containers for domain name records, with support for DNSSEC and SOA configuration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"contact": schema.StringAttribute{
				Description: "email address of the administrator responsible for this zone",
				Computed:    true,
			},
			"dnssec_enabled": schema.BoolAttribute{
				Description: "describe dnssec status\ntrue means dnssec is enabled for the zone\nfalse means dnssec is disabled for the zone",
				Computed:    true,
			},
			"dnssec_status": schema.StringAttribute{
				Description: "`dnssec_status` is the four-state lifecycle status of DNSSEC for the zone, driven by the\nparent-DS scan against the registrar. One of: pending, active, pending-disabled, disabled.\nEmpty when DNSSEC has never been enabled for the zone.",
				Computed:    true,
			},
			"dnssec_status_modified_on": schema.StringAttribute{
				Description: "`dnssec_status_modified_on` is the RFC3339 timestamp of the last `dnssec_status` change.",
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Computed: true,
			},
			"expiry": schema.Int64Attribute{
				Description: "number of seconds after which secondary name servers should stop answering request for this zone",
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
			"find_one_by": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"id": schema.ListAttribute{
						Description: "to pass several ids `id=1&id=3&id=5...`",
						Optional:    true,
						ElementType: types.Int64Type,
					},
					"case_sensitive": schema.BoolAttribute{
						Optional: true,
					},
					"client_id": schema.ListAttribute{
						Description: "to pass several `client_ids` `client_id=1&client_id=3&client_id=5...`",
						Optional:    true,
						ElementType: types.Int64Type,
					},
					"dynamic": schema.BoolAttribute{
						Description: "Zones with dynamic RRsets",
						Optional:    true,
					},
					"enabled": schema.BoolAttribute{
						Optional: true,
					},
					"exact_match": schema.BoolAttribute{
						Optional: true,
					},
					"healthcheck": schema.BoolAttribute{
						Description: "Zones with RRsets that have healthchecks",
						Optional:    true,
					},
					"iam_reseller_id": schema.ListAttribute{
						Optional:    true,
						ElementType: types.Int64Type,
					},
					"name": schema.ListAttribute{
						Description: "to pass several names `name=first&name=second...`",
						Optional:    true,
						ElementType: types.StringType,
					},
					"order_by": schema.StringAttribute{
						Description: "Field name to sort by",
						Optional:    true,
					},
					"order_direction": schema.StringAttribute{
						Description: "Ascending or descending order\nAvailable values: \"asc\", \"desc\".",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive("asc", "desc"),
						},
					},
					"reseller_id": schema.ListAttribute{
						Optional:    true,
						ElementType: types.Int64Type,
					},
					"status": schema.StringAttribute{
						Optional: true,
					},
					"updated_at_from": schema.StringAttribute{
						Optional:   true,
						CustomType: timetypes.RFC3339Type{},
					},
					"updated_at_to": schema.StringAttribute{
						Optional:   true,
						CustomType: timetypes.RFC3339Type{},
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
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(path.MatchRoot("name"), path.MatchRoot("find_one_by")),
	}
}
