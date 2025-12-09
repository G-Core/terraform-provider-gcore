// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package dns_zone

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ resource.ResourceWithConfigValidators = (*DNSZoneResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Description: "name of DNS zone",
				Required:    true,
			},
			"contact": schema.StringAttribute{
				Description: "email address of the administrator responsible for this zone",
				Optional:    true,
			},
			"expiry": schema.Int64Attribute{
				Description: "number of seconds after which secondary name servers should stop answering request for this zone",
				Optional:    true,
			},
			"nx_ttl": schema.Int64Attribute{
				Description: "Time To Live of cache",
				Optional:    true,
			},
			"primary_server": schema.StringAttribute{
				Description: "primary master name server for zone",
				Optional:    true,
			},
			"refresh": schema.Int64Attribute{
				Description: "number of seconds after which secondary name servers should query the master for the SOA record, to detect zone changes.",
				Optional:    true,
			},
			"retry": schema.Int64Attribute{
				Description: "number of seconds after which secondary name servers should retry to request the serial number",
				Optional:    true,
			},
			"serial": schema.Int64Attribute{
				Description: "Serial number for this zone or Timestamp of zone modification moment.\nIf a secondary name server slaved to this one observes an increase in this number,\nthe slave will assume that the zone has been updated and initiate a zone transfer.",
				Optional:    true,
			},
			"meta": schema.MapAttribute{
				Description: "arbitrarily data of zone in json format\nyou can specify `webhook` url and `webhook_method` here\nwebhook will get a map with three arrays: for created, updated and deleted rrsets\n`webhook_method` can be omitted, POST will be used by default",
				Optional:    true,
				ElementType: jsontypes.NormalizedType{},
			},
			"enabled": schema.BoolAttribute{
				Description: "If a zone is disabled, then its records will not be resolved on dns servers",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(true),
			},
			"warnings": schema.ListAttribute{
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
			"zone": schema.SingleNestedAttribute{
				Description: "OutputZone",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[DNSZoneZoneModel](ctx),
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Description: "ID of zone.\nThis field usually is omitted in response and available only in\ncase of getting deleted zones by admin.",
						Computed:    true,
					},
					"client_id": schema.Int64Attribute{
						Computed: true,
					},
					"contact": schema.StringAttribute{
						Description: "email address of the administrator responsible for this zone",
						Computed:    true,
					},
					"dnssec_enabled": schema.BoolAttribute{
						Description: "describe dnssec status\ntrue means dnssec is enabled for the zone\nfalse means dnssec is disabled for the zone",
						Computed:    true,
					},
					"expiry": schema.Int64Attribute{
						Description: "number of seconds after which secondary name servers should stop answering request for this zone",
						Computed:    true,
					},
					"meta": schema.StringAttribute{
						Description: "arbitrarily data of zone in json format",
						Computed:    true,
						CustomType:  jsontypes.NormalizedType{},
					},
					"name": schema.StringAttribute{
						Description: "name of DNS zone",
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
					"records": schema.ListNestedAttribute{
						Computed:   true,
						CustomType: customfield.NewNestedObjectListType[DNSZoneZoneRecordsModel](ctx),
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
					"refresh": schema.Int64Attribute{
						Description: "number of seconds after which secondary name servers should query the master for the SOA record, to detect zone changes.",
						Computed:    true,
					},
					"retry": schema.Int64Attribute{
						Description: "number of seconds after which secondary name servers should retry to request the serial number",
						Computed:    true,
					},
					"rrsets_amount": schema.SingleNestedAttribute{
						Computed:   true,
						CustomType: customfield.NewNestedObjectType[DNSZoneZoneRrsetsAmountModel](ctx),
						Attributes: map[string]schema.Attribute{
							"dynamic": schema.SingleNestedAttribute{
								Description: "Amount of dynamic RRsets in zone",
								Computed:    true,
								CustomType:  customfield.NewNestedObjectType[DNSZoneZoneRrsetsAmountDynamicModel](ctx),
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
					"serial": schema.Int64Attribute{
						Description: "Serial number for this zone or Timestamp of zone modification moment.\nIf a secondary name server slaved to this one observes an increase in this number,\nthe slave will assume that the zone has been updated and initiate a zone transfer.",
						Computed:    true,
					},
					"status": schema.StringAttribute{
						Computed: true,
					},
				},
			},
		},
	}
}

func (r *DNSZoneResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *DNSZoneResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
