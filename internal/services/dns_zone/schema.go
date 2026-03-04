// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package dns_zone

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigValidators = (*DNSZoneResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:   "name of DNS zone",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown(), stringplanmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Description:   "name of DNS zone",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown(), stringplanmodifier.RequiresReplace()},
			},
			"contact": schema.StringAttribute{
				Description:   "email address of the administrator responsible for this zone",
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"expiry": schema.Int64Attribute{
				Description:   "number of seconds after which secondary name servers should stop answering request for this zone",
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"nx_ttl": schema.Int64Attribute{
				Description:   "Time To Live of cache",
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"primary_server": schema.StringAttribute{
				Description:   "primary master name server for zone",
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"refresh": schema.Int64Attribute{
				Description:   "number of seconds after which secondary name servers should query the master for the SOA record, to detect zone changes.",
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"retry": schema.Int64Attribute{
				Description:   "number of seconds after which secondary name servers should retry to request the serial number",
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"serial": schema.Int64Attribute{
				Description: "Serial number for this zone. Server-managed: always derived from the zone's last modification timestamp.",
				Computed:    true,
			},
			"meta": schema.MapAttribute{
				Description: "arbitrarily data of zone in json format\nyou can specify `webhook` url and `webhook_method` here\nwebhook will get a map with three arrays: for created, updated and deleted rrsets\n`webhook_method` can be omitted, POST will be used by default",
				Optional:    true,
				ElementType: customfield.MetaStringType{},
			},
			"enabled": schema.BoolAttribute{
				Description: "If a zone is disabled, then its records will not be resolved on dns servers",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(true),
			},
			"dnssec_enabled": schema.BoolAttribute{
				Description:   "Whether DNSSEC is enabled for the zone. Set to true to enable DNSSEC, false to disable. Managed via a separate API endpoint.",
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			},
			"status": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"warnings": schema.ListAttribute{
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
			"records": schema.ListNestedAttribute{
				Computed:   true,
				CustomType: customfield.NewNestedObjectListType[DNSZoneRecordsModel](ctx),
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
				Computed:      true,
				CustomType:    customfield.NewNestedObjectType[DNSZoneRrsetsAmountModel](ctx),
				PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
				Attributes: map[string]schema.Attribute{
					"dynamic": schema.SingleNestedAttribute{
						Description: "Amount of dynamic RRsets in zone",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[DNSZoneRrsetsAmountDynamicModel](ctx),
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

func (r *DNSZoneResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *DNSZoneResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
