// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package dns_zone_rrset

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ resource.ResourceWithConfigValidators = (*DNSZoneRrsetResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "DNS resource record sets (RRsets) define individual DNS records such as A, AAAA, CNAME, MX, and TXT with TTL and geo-balancing settings.",
		Attributes: map[string]schema.Attribute{
			"rrset_name": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"rrset_type": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"zone_name": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"resource_records": schema.ListNestedAttribute{
				Description: "List of resource record from rrset",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "ID of the resource record",
							Computed:    true,
						},
						"content": schema.ListAttribute{
							Description: "Content of resource record\nThe exact length of the array depends on the type of rrset,\neach individual record parameter must be a separate element of the array. For example\nSRV-record: `[100, 1, 5061, \"example.com\"]`\nCNAME-record: `[ \"the.target.domain\" ]`\nA-record: `[ \"1.2.3.4\", \"5.6.7.8\" ]`\nAAAA-record: `[ \"2001:db8::1\", \"2001:db8::2\" ]`\nMX-record: `[ \"mail1.example.com\", \"mail2.example.com\" ]`\nSVCB/HTTPS-record: `[ 1, \".\", [\"alpn\", \"h3\", \"h2\"], [ \"port\", 1443 ], [ \"ipv4hint\", \"10.0.0.1\" ], [ \"ech\", \"AEn+DQBFKwAgACABWIHUGj4u+PIggYXcR5JF0gYk3dCRioBW8uJq9H4mKAAIAAEAAQABAANAEnB1YmxpYy50bHMtZWNoLmRldgAA\" ] ]`",
							Required:    true,
							ElementType: jsontypes.NormalizedType{},
						},
						"enabled": schema.BoolAttribute{
							Computed: true,
							Optional: true,
							Default:  booldefault.StaticBool(true),
						},
						"meta": schema.MapAttribute{
							Description: "This meta will be used to decide which resource record should pass\nthrough filters from the filter set",
							Optional:    true,
							ElementType: customfield.MetaStringType{},
						},
					},
				},
			},
			"ttl": schema.Int64Attribute{
				Optional: true,
			},
			"meta": schema.MapAttribute{
				Description: "Meta information for rrset",
				Optional:    true,
				ElementType: customfield.MetaStringType{},
			},
			"pickers": schema.ListNestedAttribute{
				Description: "Set of pickers",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Description: "Filter type\nAvailable values: \"geodns\", \"asn\", \"country\", \"continent\", \"region\", \"ip\", \"geodistance\", \"weighted_shuffle\", \"default\", \"first_n\".",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"geodns",
									"asn",
									"country",
									"continent",
									"region",
									"ip",
									"geodistance",
									"weighted_shuffle",
									"default",
									"first_n",
								),
							},
						},
						"limit": schema.Int64Attribute{
							Description: "Limits the number of records returned by the filter\nCan be a positive value for a specific limit. Use zero or leave it blank to indicate no limits.",
							Computed:    true,
							Optional:    true,
						},
						"strict": schema.BoolAttribute{
							Description: "if strict=false, then the filter will return all records if no records match the filter",
							Computed:    true,
							Optional:    true,
						},
					},
				},
			},
			"filter_set_id": schema.Int64Attribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"type": schema.StringAttribute{
				Description:   "RRSet type\nAvailable values: \"A\", \"AAAA\", \"NS\", \"CNAME\", \"MX\", \"TXT\", \"SRV\", \"SOA\".",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"A",
						"AAAA",
						"NS",
						"CNAME",
						"MX",
						"TXT",
						"SRV",
						"SOA",
					),
				},
			},
			"warning": schema.StringAttribute{
				Description: "Warning about some possible side effects without strictly disallowing operations on rrset\nreadonly\nDeprecated: use Warnings instead",
				Computed:    true,
			},
			"warnings": schema.ListNestedAttribute{
				Description: "Warning about some possible side effects without strictly disallowing operations on rrset\nreadonly",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[DNSZoneRrsetWarningsModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Computed: true,
						},
						"message": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (r *DNSZoneRrsetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *DNSZoneRrsetResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
