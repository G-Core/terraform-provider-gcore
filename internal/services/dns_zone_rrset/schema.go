// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package dns_zone_rrset

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ resource.ResourceWithConfigValidators = (*DNSZoneRrsetResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
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
						"content": schema.ListAttribute{
							Description: "Content of resource record\nThe exact length of the array depends on the type of rrset,\neach individual record parameter must be a separate element of the array. For example\n- SRV-record: `[100, 1, 5061, \"example.com\"]`\n- CNAME-record: `[ \"the.target.domain\" ]`\n- A-record: `[ \"1.2.3.4\", \"5.6.7.8\" ]`\n- AAAA-record: `[ \"2001:db8::1\", \"2001:db8::2\" ]`\n- MX-record: `[ \"mail1.example.com\", \"mail2.example.com\" ]`\n- SVCB/HTTPS-record: `[ 1, \".\", [\"alpn\", \"h3\", \"h2\"], [ \"port\", 1443 ], [ \"ipv4hint\", \"10.0.0.1\" ], [ \"ech\", \"AEn+DQBFKwAgACABWIHUGj4u+PIggYXcR5JF0gYk3dCRioBW8uJq9H4mKAAIAAEAAQABAANAEnB1YmxpYy50bHMtZWNoLmRldgAA\" ] ]`",
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
							ElementType: jsontypes.NormalizedType{},
						},
					},
				},
				PlanModifiers: []planmodifier.List{listplanmodifier.RequiresReplace()},
			},
			"ttl": schema.Int64Attribute{
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"meta": schema.MapAttribute{
				Description:   "Meta information for rrset",
				Optional:      true,
				ElementType:   jsontypes.NormalizedType{},
				PlanModifiers: []planmodifier.Map{mapplanmodifier.RequiresReplace()},
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
							Optional:    true,
						},
						"strict": schema.BoolAttribute{
							Description: "if strict=false, then the filter will return all records if no records match the filter",
							Optional:    true,
						},
					},
				},
				PlanModifiers: []planmodifier.List{listplanmodifier.RequiresReplace()},
			},
			"filter_set_id": schema.Int64Attribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"type": schema.StringAttribute{
				Description: "RRSet type\nAvailable values: \"A\", \"AAAA\", \"NS\", \"CNAME\", \"MX\", \"TXT\", \"SRV\", \"SOA\".",
				Computed:    true,
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
			"updated_at": schema.StringAttribute{
				Description: "Timestamp marshals/unmarshals date and time as timestamp in json",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
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
